package service_test

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestAddTransaction(t *testing.T) {
	//create fake account
	account := createRandomAccount(t)
	// create random investment
	investment := createRandomInvestmnet(t, account.ID)
	// create transaction params
	trans_params := []dtos.CreateTransactionDto{
		{
			AccountId:    account.ID,
			InvestmentId: investment.ID,
			Ticker:       investment.Ticker,
			TradingDate: pgtype.Timestamp{
				Time:  time.Date(2024, 07, 31, 9, 0, 0, 0, time.UTC),
				Valid: true,
			},
			Trade:       db.TradeTypeBUY,
			Volume:      1000,
			OrderPrice:  33500,
			MatchVolume: 1000,
			MatchPrice:  33500,
			Fee:         50250,
			Tax:         0,
			Status:      db.TransactionStatusCOMPLETED,
		},
		{
			AccountId:    account.ID,
			InvestmentId: investment.ID,
			Ticker:       investment.Ticker,
			TradingDate: pgtype.Timestamp{
				Time:  time.Date(2024, 8, 2, 9, 0, 0, 0, time.UTC),
				Valid: true,
			},
			Trade:       db.TradeTypeSELL,
			Volume:      100,
			OrderPrice:  31400,
			MatchVolume: 100,
			MatchPrice:  31400,
			Fee:         4710,
			Tax:         3140,
			Status:      db.TransactionStatusCOMPLETED,
		},
		{
			AccountId:    account.ID,
			InvestmentId: investment.ID,
			Ticker:       investment.Ticker,
			TradingDate: pgtype.Timestamp{
				Time:  time.Date(2024, 8, 3, 10, 0, 0, 0, time.UTC),
				Valid: true,
			},
			Trade:       db.TradeTypeSELL,
			Volume:      100,
			OrderPrice:  31400,
			MatchVolume: 100,
			MatchPrice:  31400,
			Fee:         4710,
			Tax:         3140,
			Status:      db.TransactionStatusCOMPLETED,
		},
	}
	for _, param := range trans_params {
		transaction, err := tranService.AddTransaction(context.Background(), param)
		require.NoError(t, err)
		require.NotEmpty(t, transaction)
		require.Equal(t, transaction.InvestmentID, investment.ID)
		require.Equal(t, transaction.Ticker, param.Ticker)
		require.Equal(t, transaction.TradingDate, param.TradingDate)
		require.Equal(t, transaction.Trade, param.Trade)
		require.Equal(t, transaction.Volume, param.Volume)
		require.Equal(t, transaction.OrderPrice, param.OrderPrice)
		require.Equal(t, transaction.MatchVolume, param.MatchVolume)
		require.Equal(t, transaction.Fee, param.Fee)
		require.Equal(t, transaction.Tax, param.Tax)
		require.Equal(t, transaction.Status, param.Status)
		require.Equal(t, param.Trade, transaction.Trade)
		require.Equal(t, transaction.MatchValue, transaction.MatchPrice*transaction.Volume)

		// check entries
		// ....
		// check account's balance
		accountdb, err := accService.GetById(context.Background(), account.ID)
		if transaction.Trade == db.TradeTypeBUY {
			account.Balance -= (transaction.MatchValue + transaction.Fee + transaction.Tax)
		} else {
			account.Balance += (transaction.MatchValue - transaction.Fee - transaction.Tax)
		}
		require.NoError(t, err)
		require.Equal(t, accountdb.Balance, account.Balance)

		// check with investment
		dbInvestment, err := investmentService.GetById(context.Background(), param.InvestmentId)
		require.NoError(t, err)

		require.Equal(t, investment.Ticker, dbInvestment.Ticker)
		require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), dbInvestment.UpdatedDate.Time.Truncate(24*time.Hour))
		require.Equal(t, transaction.Cost, dbInvestment.CapitalCost)
		require.Equal(t, transaction.CostOfGoodsSold, dbInvestment.CapitalCost*transaction.MatchVolume)

		if transaction.Trade == db.TradeTypeBUY {
			investment.BuyValue += transaction.MatchValue
			investment.BuyVolume += transaction.MatchVolume
			investment.CurrentVolume += transaction.MatchVolume

			require.Equal(t, investment.BuyValue, dbInvestment.BuyValue)
			require.Equal(t, investment.BuyVolume, dbInvestment.BuyVolume)

			require.Equal(t, transaction.Return, int64(0))
		} else if transaction.Trade == db.TradeTypeSELL {
			investment.SellValue += transaction.MatchValue
			investment.SellVolume += transaction.MatchVolume
			investment.CurrentVolume -= transaction.MatchVolume

			require.Equal(t, investment.SellValue, dbInvestment.SellValue)
			require.Equal(t, investment.SellVolume, dbInvestment.SellVolume)

			returnValue := ((param.MatchPrice - dbInvestment.CapitalCost) * param.MatchVolume) - param.Fee - param.Tax
			require.Equal(t, transaction.Return, returnValue)
		}
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
}
