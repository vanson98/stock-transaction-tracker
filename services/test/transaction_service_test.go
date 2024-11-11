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

	trans_params := []dtos.CreateTransactionDto{
		{
			AccountId:    account.ID,
			InvestmentId: investment.ID,
			Ticker:       investment.Ticker,
			TradingDate: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
			Trade:       db.TradeTypeBUY,
			Volume:      100,
			OrderPrice:  34750,
			MatchVolume: 100,
			MatchPrice:  34750,
			Fee:         5212,
			Tax:         0,
			Status:      db.TransactionStatusCOMPLETED,
		},
		{
			AccountId:    account.ID,
			InvestmentId: investment.ID,
			Ticker:       investment.Ticker,
			TradingDate: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
			Trade:       db.TradeTypeBUY,
			Volume:      100,
			OrderPrice:  35350,
			MatchVolume: 100,
			MatchPrice:  35350,
			Fee:         5302,
			Tax:         0,
			Status:      db.TransactionStatusCOMPLETED,
		},
	}
	var totalBuyValue int64
	for _, param := range trans_params {
		transaction, err := tranService.AddTransaction(context.Background(), param)
		require.NoError(t, err)
		require.NotEmpty(t, transaction)
		require.Equal(t, transaction.InvestmentID, investment.ID)
		require.Equal(t, transaction.Ticker, param.Ticker)
		//require.Equal(t, transaction.TradingDate.Time, param.TradingDate.Time)
		require.Equal(t, transaction.Trade, param.Trade)
		require.Equal(t, transaction.Volume, param.Volume)
		require.Equal(t, transaction.OrderPrice, param.OrderPrice)
		require.Equal(t, transaction.MatchVolume, param.MatchVolume)
		require.Equal(t, transaction.Fee, param.Fee)
		require.Equal(t, transaction.Tax, param.Tax)
		require.Equal(t, transaction.Status, param.Status)
		require.Equal(t, transaction.Return, int64(0))
		// check entries
		// ....
		// check account
		newBalance := account.Balance - transaction.MatchValue - transaction.Fee - transaction.Tax
		account, err = accService.GetById(context.Background(), account.ID)
		require.NoError(t, err)
		require.Equal(t, account.Balance, newBalance)

		// check capital cost
		investment, err = investmentService.GetById(context.Background(), param.InvestmentId)
		require.NoError(t, err)

		require.Equal(t, transaction.MatchValue, transaction.MatchPrice*transaction.Volume)
		require.Equal(t, transaction.Cost, investment.CapitalCost)
		require.Equal(t, transaction.CostOfGoodsSold, investment.CapitalCost*transaction.MatchVolume)
		// check match value
		if transaction.Trade == db.TradeTypeBUY {
			totalBuyValue += transaction.MatchValue
			require.Equal(t, investment.BuyValue, totalBuyValue)
		}

	}
}
