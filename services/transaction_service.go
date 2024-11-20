package services

import (
	"context"
	"fmt"
	"math"
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
	return &transactionService{
		store: store,
	}
}

// GetPaging implements sv_interface.ITransactionService.
func (t *transactionService) GetPaging(ctx context.Context, ticker string) ([]db.Transaction, error) {
	return t.store.GetTransactionsPaging(ctx, ticker)
}

// GetById implements sv_interface.ITransactionService.
func (t *transactionService) GetById(ctx context.Context, id int64) (db.Transaction, error) {
	return t.store.GetTransactionById(ctx, id)
}

// CreateNew implements sv_interface.ITransactionService.
func (t *transactionService) AddTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error) {
	if arg.Trade == db.TradeTypeBUY {
		return t.insertBuyingTransaction(ctx, arg)
	} else {
		return t.insertSellingTransaction(ctx, arg)
	}
}

func (t *transactionService) insertBuyingTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error) {
	result, err := t.store.ExecTx(ctx, func(q *db.Queries) (interface{}, error) {
		// check account balance
		account, err := t.store.GetAccountById(ctx, arg.AccountId)
		if err != nil {
			return nil, err
		}
		totalTransactionValue := (arg.MatchPrice * arg.MatchVolume) + arg.Fee + arg.Tax
		if account.Balance < totalTransactionValue {
			return nil, fmt.Errorf("account balance is less than transation cost")
		}

		// create entry
		entry, err := t.store.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: arg.AccountId,
			Amount:    -totalTransactionValue,
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

		// calculate capital cost for each shares
		investment, err := t.store.GetInvestmentById(ctx, arg.InvestmentId)
		if err != nil {
			return nil, err
		}
		currentInvestmentValue := investment.CurrentVolume * investment.CapitalCost
		roundUpCapitalCost := math.Round(
			(float64(currentInvestmentValue) + float64(totalTransactionValue)) /
				(float64(investment.CurrentVolume) + float64(arg.MatchVolume)))
		investment.CapitalCost = int64(roundUpCapitalCost)

		// create transaction
		transaction, err := t.store.CreateTransaction(ctx, db.CreateTransactionParams{
			InvestmentID:    arg.InvestmentId,
			Ticker:          arg.Ticker,
			TradingDate:     arg.TradingDate,
			Trade:           arg.Trade,
			Volume:          arg.Volume,
			OrderPrice:      arg.OrderPrice,
			MatchVolume:     arg.MatchVolume,
			MatchPrice:      arg.MatchPrice,
			MatchValue:      arg.MatchPrice * arg.MatchVolume,
			Fee:             arg.Fee,
			Tax:             arg.Tax,
			Return:          0,
			Status:          arg.Status,
			Cost:            investment.CapitalCost,
			CostOfGoodsSold: investment.CapitalCost * arg.MatchVolume,
		})
		if err != nil {
			return nil, err
		}

		// update investment
		investment.BuyVolume += transaction.MatchVolume
		investment.BuyValue += transaction.MatchValue
		investment.Fee += transaction.Fee
		investment.Tax += transaction.Tax
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
		return transaction, nil
	})
	if err != nil {
		return db.Transaction{}, err
	}
	transaction, ok := result.(db.Transaction)
	if !ok {
		err = fmt.Errorf("can not convert db tx result to transaction type")
		return db.Transaction{}, err
	}
	return transaction, err
}

func (t *transactionService) insertSellingTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error) {
	result, err := t.store.ExecTx(ctx, func(query *db.Queries) (interface{}, error) {
		// get investment
		investment, err := query.GetInvestmentById(ctx, arg.InvestmentId)
		if err != nil {
			return db.Transaction{}, err
		}

		// create a transaction
		transaction, err := t.store.CreateTransaction(ctx, db.CreateTransactionParams{
			InvestmentID:    arg.InvestmentId,
			Ticker:          arg.Ticker,
			TradingDate:     arg.TradingDate,
			Trade:           arg.Trade,
			Volume:          arg.Volume,
			OrderPrice:      arg.OrderPrice,
			MatchVolume:     arg.MatchVolume,
			MatchPrice:      arg.MatchPrice,
			MatchValue:      arg.MatchPrice * arg.MatchVolume,
			Fee:             arg.Fee,
			Tax:             arg.Tax,
			Cost:            investment.CapitalCost,
			CostOfGoodsSold: investment.CapitalCost * arg.MatchVolume,
			Return:          (arg.MatchPrice * arg.MatchVolume) - arg.Fee - arg.Tax - (investment.CapitalCost * arg.MatchVolume),
			Status:          arg.Status,
		})
		if err != nil {
			return db.Transaction{}, err
		}

		// create account's entry
		entry, err := query.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: arg.AccountId,
			Amount:    transaction.MatchValue - transaction.Fee - transaction.Tax,
		})
		if err != nil {
			return db.Transaction{}, err
		}

		// update account' balance
		query.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			Amount: entry.Amount,
			ID:     arg.AccountId,
		})

		// update investment
		err = query.UpdateInvestmentWhenSeling(ctx, db.UpdateInvestmentWhenSelingParams{
			ID:                    transaction.InvestmentID,
			SellTransactionVolume: transaction.MatchVolume,
			SellTransactionValue:  transaction.MatchValue,
			TransactionFee:        transaction.Fee,
			TransactionTax:        transaction.Tax,
			UpdatedDate:           arg.TradingDate,
		})
		if err != nil {
			return db.Transaction{}, err
		}

		return transaction, nil
	})
	if err != nil {
		return db.Transaction{}, nil
	}
	transaction, ok := result.(db.Transaction)
	if !ok {
		err = fmt.Errorf("can not convert db tx result to transaction type")
		return db.Transaction{}, err
	}
	return transaction, err
}
