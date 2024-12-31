package service_test

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
	"stt/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomTransaction(t *testing.T, investmentId int64, ticker string, trade db.TradeType, matchVolume int64) db.Transaction {
	// create transaction params
	param := dtos.CreateTransactionDto{
		InvestmentId: investmentId,
		TradingDate: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		Trade:       trade,
		Volume:      util.RandomInt(100, 5000),
		OrderPrice:  util.RandomInt(10000, 100000),
		MatchVolume: matchVolume,
		MatchPrice:  util.RandomInt(10000, 100000),
		Fee:         util.RandomInt(10000, 20000),
		Tax:         util.RandomInt(10000, 20000),
		Status:      db.TransactionStatusCOMPLETED,
	}

	transaction, err := tranService.CreateTransaction(context.Background(), param)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)
	require.Equal(t, transaction.InvestmentID, investmentId)
	require.Equal(t, transaction.Ticker, ticker)
	require.Equal(t, transaction.TradingDate.Time.Truncate(24*time.Hour), param.TradingDate.Time.UTC().Truncate(24*time.Hour))
	require.Equal(t, transaction.Trade, param.Trade)
	require.Equal(t, transaction.Volume, param.Volume)
	require.Equal(t, transaction.OrderPrice, param.OrderPrice)
	require.Equal(t, transaction.MatchVolume, param.MatchVolume)
	require.Equal(t, transaction.Fee, param.Fee)
	require.Equal(t, transaction.Tax, param.Tax)
	require.Equal(t, transaction.Status, param.Status)
	require.Equal(t, param.Trade, transaction.Trade)
	require.Equal(t, transaction.MatchValue, param.MatchPrice*param.MatchVolume)

	return transaction
}

func TestCreateCircleTransaction(t *testing.T) {
	user := createRandomUser(t)
	acc := createRandomAccount(t, user.Username)
	investment := createDefaultInvestmnet(t, acc.ID)
	buyVolume := util.RandomInt(500, 1000)
	createRandomBuyTransaction(t, acc.ID, investment.ID, buyVolume)
	sellVolume := util.RandomInt(100, 500)
	createRandomSellTransaction(t, acc.ID, investment.ID, sellVolume)
}

func createRandomBuyTransaction(t *testing.T, accountId int64, investmentId int64, matchVolume int64) {
	investment, err := investmentService.GetById(context.Background(), investmentId)
	require.NoError(t, err)
	account, err := accService.GetById(context.Background(), accountId)
	require.NoError(t, err)
	transaction := createRandomTransaction(t, investment.ID, investment.Ticker, db.TradeTypeBUY, matchVolume)

	// check account's balance
	updatedAccount, err := accService.GetById(context.Background(), account.ID)
	account.Balance -= (transaction.MatchValue + transaction.Fee + transaction.Tax)
	require.NoError(t, err)
	require.Equal(t, updatedAccount.Balance, account.Balance)

	// check investment
	dbInvestment, err := investmentService.GetById(context.Background(), investment.ID)
	require.NoError(t, err)
	require.Equal(t, investment.Ticker, dbInvestment.Ticker)
	require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), dbInvestment.UpdatedDate.Time.Truncate(24*time.Hour))
	require.Equal(t, transaction.Cost, dbInvestment.CapitalCost)
	require.Equal(t, transaction.CostOfGoodsSold, dbInvestment.CapitalCost*transaction.MatchVolume)

	investment.BuyValue += transaction.MatchValue
	investment.BuyVolume += transaction.MatchVolume
	investment.CurrentVolume += transaction.MatchVolume
	require.Equal(t, investment.BuyValue, dbInvestment.BuyValue)
	require.Equal(t, investment.BuyVolume, dbInvestment.BuyVolume)
	require.Equal(t, transaction.Return, int64(0))

	investment.Fee += transaction.Fee
	investment.Tax += transaction.Tax
	require.Equal(t, investment.CurrentVolume, dbInvestment.CurrentVolume)
	require.Equal(t, investment.Fee, dbInvestment.Fee)
	require.Equal(t, investment.Tax, dbInvestment.Tax)
	require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), dbInvestment.UpdatedDate.Time.Truncate(24*time.Hour))

	if dbInvestment.CurrentVolume > 0 {
		require.Equal(t, db.InvestmentStatusActive, dbInvestment.Status)
	}

}

func createRandomSellTransaction(t *testing.T, accountId int64, investmentId int64, matchVolume int64) {
	investment, err := investmentService.GetById(context.Background(), investmentId)
	require.NoError(t, err)
	account, err := accService.GetById(context.Background(), accountId)
	require.NoError(t, err)
	transaction := createRandomTransaction(t, investment.ID, investment.Ticker, db.TradeTypeSELL, matchVolume)

	// check account's balance
	updatedAccount, err := accService.GetById(context.Background(), account.ID)
	account.Balance += (transaction.MatchValue - transaction.Fee - transaction.Tax)
	require.NoError(t, err)
	require.Equal(t, updatedAccount.Balance, account.Balance)

	// check investment
	updatedInvestment, err := investmentService.GetById(context.Background(), investment.ID)
	require.NoError(t, err)
	require.Equal(t, investment.Ticker, updatedInvestment.Ticker)
	require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), updatedInvestment.UpdatedDate.Time.Truncate(24*time.Hour))
	require.Equal(t, transaction.Cost, updatedInvestment.CapitalCost)
	require.Equal(t, transaction.CostOfGoodsSold, updatedInvestment.CapitalCost*transaction.MatchVolume)

	investment.SellValue += transaction.MatchValue
	investment.SellVolume += transaction.MatchVolume
	investment.CurrentVolume -= transaction.MatchVolume
	require.Equal(t, investment.SellValue, updatedInvestment.SellValue)
	require.Equal(t, investment.SellVolume, updatedInvestment.SellVolume)
	returnValue := ((transaction.MatchPrice - updatedInvestment.CapitalCost) * transaction.MatchVolume) - transaction.Fee - transaction.Tax
	require.Equal(t, transaction.Return, returnValue)

	investment.Fee += transaction.Fee
	investment.Tax += transaction.Tax
	require.Equal(t, investment.CurrentVolume, updatedInvestment.CurrentVolume)
	require.Equal(t, investment.Fee, updatedInvestment.Fee)
	require.Equal(t, investment.Tax, updatedInvestment.Tax)
	require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), updatedInvestment.UpdatedDate.Time.Truncate(24*time.Hour))

	if updatedInvestment.CurrentVolume == 0 {
		require.Equal(t, db.InvestmentStatusSellout, updatedInvestment.Status)
	}
}
