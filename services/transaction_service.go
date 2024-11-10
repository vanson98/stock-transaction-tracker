package services

import (
	"context"
	"fmt"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
	sv_interface "stt/services/interfaces"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type transactionService struct {
	store db.IStore
}

func InitTransactionService(store db.IStore) sv_interface.ITransactionService {
	return transactionService{
		store: store,
	}
}

// GetById implements sv_interface.ITransactionService.
func (t transactionService) GetById(ctx context.Context, id int64) (db.Transaction, error) {
	return t.store.GetTransactionById(ctx, id)
}

// CreateNew implements sv_interface.ITransactionService.
func (t transactionService) CreateBuyingTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error) {
	result, err := t.store.ExecTx(ctx, func(q *db.Queries) (interface{}, error) {
		// check account balance
		account, err := t.store.GetAccountById(ctx, arg.AccountId)
		if err != nil {
			return nil, err
		}
		transctionTempCost := (arg.MatchPrice * arg.MatchVolume) + arg.Fee + arg.Tax
		if account.Balance < transctionTempCost {
			return nil, fmt.Errorf("account balance is less than transation cost")
		}

		// create entry
		entry, err := t.store.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: arg.AccountId,
			Amount:    -transctionTempCost,
			Type:      db.EntryTypeIT,
		})
		if err != nil {
			return nil, err
		}

		//update account balance
		account, err = q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			ID:     account.ID,
			Amount: entry.Amount,
		})
		if err != nil {
			return nil, err
		}

		// create transaction
		transaction, err := t.store.CreateTransaction(ctx, db.CreateTransactionParams{
			InvestmentID: arg.InvestmentId,
			Ticker:       arg.Ticker,
			TradingDate:  arg.TradingDate,
			Trade:        arg.Trade,
			Volume:       arg.Volume,
			OrderPrice:   arg.OrderPrice,
			MatchVolume:  arg.MatchVolume,
			MatchPrice:   arg.MatchPrice,
			MatchValue:   arg.MatchPrice * arg.MatchVolume,
			Fee:          arg.Fee,
			Tax:          arg.Tax,
			Return:       0,
			Status:       arg.Status,
		})
		if err != nil {
			return nil, err
		}

		// update investment
		investment, err := t.store.GetInvestmentById(ctx, arg.InvestmentId)
		if err != nil {
			return nil, err
		}
		investment.BuyVolume += transaction.MatchVolume
		investment.BuyValue += transaction.MatchValue
		investment.Fee += transaction.Fee
		investment.Tax += transaction.Tax
		currentInvestmentCgs := investment.CurrentVolume * investment.CapitalCost
		investment.CapitalCost = (currentInvestmentCgs + transctionTempCost) / (investment.CurrentVolume + transaction.MatchVolume)
		investment.CurrentVolume += transaction.MatchVolume
		err = t.store.UpdateInvestmentWhenBuying(ctx, db.UpdateInvestmentWhenBuyingParams{
			ID:            investment.ID,
			BuyVolume:     investment.BuyVolume,
			BuyValue:      investment.BuyValue,
			CapitalCost:   investment.CapitalCost,
			CurrentVolume: investment.CurrentVolume,
			Fee:           investment.Fee,
			Tax:           investment.Tax,
			UpdatedDate: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			return nil, err
		}

		// update transaction cost
		transaction, err = t.store.UpdateTransactionCost(ctx, db.UpdateTransactionCostParams{
			ID:              transaction.ID,
			Cost:            investment.CapitalCost,
			CostOfGoodsSold: investment.CapitalCost * transaction.MatchVolume,
		})
		if err != nil {
			return nil, err
		}
		return transaction, nil
	})
	transaction, ok := result.(db.Transaction)
	if !ok {
		return db.Transaction{}, err
	}
	return transaction, err
}
