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

func createRandomTransaction(t *testing.T, investment *db.Investment, trade db.TradeType, matchVolume int64) db.Transaction {

	// create transaction params
	param := dtos.CreateTransactionDto{
		InvestmentId: investment.ID,
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
	require.Equal(t, transaction.InvestmentID, investment.ID)
	require.Equal(t, transaction.Ticker, investment.Ticker)
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

func TestCreateBuySellTransaction(t *testing.T) {
	user := createRandomUser(t)
	acc := createRandomAccount(t, user)
	investment := createDefaultInvestmnet(t, acc.ID)
	buyVolume := util.RandomInt(500, 1000)
	createRandomBuyTransaction(t, &acc, &investment, buyVolume)
	sellVolume := util.RandomInt(100, 500)
	createRamdomSellTransaction(t, &acc, &investment, sellVolume)
}

func createRandomBuyTransaction(t *testing.T, acc *db.Account, investment *db.Investment, matchVolume int64) {
	transaction := createRandomTransaction(t, investment, db.TradeTypeBUY, matchVolume)

	// check account's balance
	accountdb, err := accService.GetById(context.Background(), acc.ID)
	acc.Balance -= (transaction.MatchValue + transaction.Fee + transaction.Tax)
	require.NoError(t, err)
	require.Equal(t, accountdb.Balance, acc.Balance)

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

func createRamdomSellTransaction(t *testing.T, acc *db.Account, investment *db.Investment, matchVolume int64) {
	transaction := createRandomTransaction(t, investment, db.TradeTypeSELL, matchVolume)

	// check account's balance
	accountdb, err := accService.GetById(context.Background(), acc.ID)
	acc.Balance += (transaction.MatchValue - transaction.Fee - transaction.Tax)
	require.NoError(t, err)
	require.Equal(t, accountdb.Balance, acc.Balance)

	// check investment
	dbInvestment, err := investmentService.GetById(context.Background(), investment.ID)
	require.NoError(t, err)
	require.Equal(t, investment.Ticker, dbInvestment.Ticker)
	require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), dbInvestment.UpdatedDate.Time.Truncate(24*time.Hour))
	require.Equal(t, transaction.Cost, dbInvestment.CapitalCost)
	require.Equal(t, transaction.CostOfGoodsSold, dbInvestment.CapitalCost*transaction.MatchVolume)

	investment.SellValue += transaction.MatchValue
	investment.SellVolume += transaction.MatchVolume
	investment.CurrentVolume -= transaction.MatchVolume
	require.Equal(t, investment.SellValue, dbInvestment.SellValue)
	require.Equal(t, investment.SellVolume, dbInvestment.SellVolume)
	returnValue := ((transaction.MatchPrice - dbInvestment.CapitalCost) * transaction.MatchVolume) - transaction.Fee - transaction.Tax
	require.Equal(t, transaction.Return, returnValue)

	investment.Fee += transaction.Fee
	investment.Tax += transaction.Tax
	require.Equal(t, investment.CurrentVolume, dbInvestment.CurrentVolume)
	require.Equal(t, investment.Fee, dbInvestment.Fee)
	require.Equal(t, investment.Tax, dbInvestment.Tax)
	require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), dbInvestment.UpdatedDate.Time.Truncate(24*time.Hour))

	if dbInvestment.CurrentVolume == 0 {
		require.Equal(t, db.InvestmentStatusSellout, dbInvestment.Status)
	}
}
