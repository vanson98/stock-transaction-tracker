package services

import (
	"context"
	"fmt"
	"math"
	"slices"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
	sv_interface "stt/services/interfaces"
	"time"

	"github.com/jackc/pgx/v5"
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

func (t *transactionService) CreateTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error) {
	if arg.Trade == db.TradeTypeBUY {
		return t.createBuyingTransaction(ctx, arg)
	} else if arg.Trade == db.TradeTypeSELL {
		return t.createSellingTransaction(ctx, arg)
	}
	return db.Transaction{}, fmt.Errorf("trading type is not valid")
}

func (t *transactionService) ImportTransactions(ctx context.Context, accountId int64, transactions []db.Transaction) ([]db.Transaction, error) {
	// Get account
	account, err := t.store.GetAccountById(ctx, accountId)
	if err != nil {
		return nil, err
	}
	slices.Reverse(transactions)
	txResult, err := t.store.ExecTx(ctx, func(q *db.Queries) (interface{}, error) {
		transactionsInserted := make([]db.Transaction, 0)
		for _, importTrans := range transactions {
			// check investment exist
			investment, err := q.GetInvestmentByTicker(ctx, db.GetInvestmentByTickerParams{
				Ticker:    importTrans.Ticker,
				AccountID: account.ID,
			})
			if err != nil {
				if err == pgx.ErrNoRows {
					investment, err = q.CreateInvestment(ctx, db.CreateInvestmentParams{
						AccountID: accountId,
						Ticker:    importTrans.Ticker,
						Status:    db.InvestmentStatusInactive,
					})
				}
				if err != nil {
					return nil, err
				}
			}

			// create transaction
			transaction, err := q.CreateTransaction(ctx, db.CreateTransactionParams{
				InvestmentID:    investment.ID,
				Ticker:          importTrans.Ticker,
				TradingDate:     importTrans.TradingDate,
				Trade:           importTrans.Trade,
				Volume:          importTrans.Volume,
				OrderPrice:      importTrans.OrderPrice,
				MatchVolume:     importTrans.MatchVolume,
				MatchPrice:      importTrans.MatchPrice,
				MatchValue:      importTrans.MatchValue,
				Fee:             importTrans.Fee,
				Tax:             importTrans.Tax,
				Return:          importTrans.Return,
				Status:          importTrans.Status,
				Cost:            importTrans.Cost,
				CostOfGoodsSold: importTrans.Cost * importTrans.MatchVolume,
			})
			if err != nil {
				return nil, err
			}

			// update investment and account balance
			if importTrans.Trade == db.TradeTypeBUY {
				_, err = t.checkBuyingTransaction(ctx, &investment, &transaction, account, q)
			} else if transaction.Trade == db.TradeTypeSELL {
				err = t.checkSellingTransaction(ctx, &investment, &transaction, q)
			}
			if err != nil {
				return nil, err
			}
			transactionsInserted = append(transactionsInserted, transaction)
		}
		return transactionsInserted, nil
	})
	if err != nil {
		return nil, err
	}
	savedTransactions, ok := txResult.([]db.Transaction)
	if !ok {
		return nil, fmt.Errorf("can't not convert db transaction result")
	}
	if len(savedTransactions) != len(transactions) {
		return nil, fmt.Errorf("all transactions are imported in DB but input amount and output amount does not match")
	}
	return transactions, nil
}

func (t *transactionService) ImportTransaction(ctx context.Context, accountId int64, transaction db.Transaction) (db.Transaction, error) {
	// Get account
	account, err := t.store.GetAccountById(ctx, accountId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return db.Transaction{}, fmt.Errorf("account not found")
		}
		return db.Transaction{}, err
	}
	txResult, err := t.store.ExecTx(ctx, func(q *db.Queries) (interface{}, error) {
		// check investment exist
		investment, err := q.GetInvestmentByTicker(ctx, db.GetInvestmentByTickerParams{
			Ticker:    transaction.Ticker,
			AccountID: account.ID,
		})
		if err != nil {
			if err == pgx.ErrNoRows {
				investment, err = q.CreateInvestment(ctx, db.CreateInvestmentParams{
					AccountID: accountId,
					Ticker:    transaction.Ticker,
					Status:    db.InvestmentStatusInactive,
				})
			}
			if err != nil {
				return db.Transaction{}, err
			}
		}
		// create transaction

		if err != nil {
			return db.Transaction{}, err
		}
		var returnError int64 = 0
		// update investment and account balance
		if transaction.Trade == db.TradeTypeBUY {
			costError, ersr := t.checkBuyingTransaction(ctx, &investment, &transaction, account, q)
			err = ersr
			returnError = costError * transaction.MatchVolume
		} else if transaction.Trade == db.TradeTypeSELL {
			returnError = (investment.CapitalCost - transaction.Cost) * transaction.MatchVolume
			if math.Abs(float64(returnError)) > 5000 {
				return db.Transaction{}, fmt.Errorf("transaction's return error too large - transaction cost :%v , investment cost: %v", transaction.Cost, investment.CapitalCost)
			}
			err = t.checkSellingTransaction(ctx, &investment, &transaction, q)
		}
		if err != nil {
			return db.Transaction{}, err
		}
		transaction, err = q.CreateTransaction(ctx, db.CreateTransactionParams{
			InvestmentID:    investment.ID,
			Ticker:          transaction.Ticker,
			TradingDate:     transaction.TradingDate,
			Trade:           transaction.Trade,
			Volume:          transaction.Volume,
			OrderPrice:      transaction.OrderPrice,
			MatchVolume:     transaction.MatchVolume,
			MatchPrice:      transaction.MatchPrice,
			MatchValue:      transaction.MatchValue,
			Fee:             transaction.Fee,
			Tax:             transaction.Tax,
			Return:          transaction.Return,
			Status:          transaction.Status,
			Cost:            transaction.Cost,
			CostOfGoodsSold: transaction.Cost * transaction.MatchVolume,
			ReturnError:     returnError,
			InsertedDate: pgtype.Timestamp{
				Time:  time.Now().UTC(),
				Valid: true,
			},
		})
		if err != nil {
			return db.Transaction{}, err
		}
		return transaction, nil
	})
	if err != nil {
		return db.Transaction{}, err
	}
	tx, ok := txResult.(db.Transaction)
	if !ok {
		return db.Transaction{}, fmt.Errorf("db transaction result's underlying value is not valid")
	}
	return tx, nil
}

func (t *transactionService) GetPaging(ctx context.Context, param db.GetTransactionsPagingParams) ([]db.GetTransactionsPagingRow, error) {
	return t.store.GetTransactionsPaging(ctx, param)
}

func (t *transactionService) GetSummarizeInfo(ctx context.Context, param db.GetTransactionSummarizeInfoParams) (db.GetTransactionSummarizeInfoRow, error) {
	return t.store.GetTransactionSummarizeInfo(ctx, param)
}

func (t *transactionService) createBuyingTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error) {
	result, err := t.store.ExecTx(ctx, func(q *db.Queries) (interface{}, error) {
		// get investment
		investment, err := t.store.GetInvestmentById(ctx, arg.InvestmentId)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, fmt.Errorf("Investment not found")
			}
			return nil, err
		}
		// check account balance
		account, err := t.store.GetAccountById(ctx, investment.AccountID)
		if err != nil {
			return nil, err
		}
		totalTransactionValue := (arg.MatchPrice * arg.MatchVolume) + arg.Fee + arg.Tax
		if account.Balance < totalTransactionValue {
			return nil, fmt.Errorf("account balance is less than transation cost")
		}

		// create entry
		entry, err := t.store.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: investment.AccountID,
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
		currentInvestmentValue := investment.CurrentVolume * investment.CapitalCost
		roundUpCapitalCost := math.Round(
			(float64(currentInvestmentValue) + float64(totalTransactionValue)) /
				(float64(investment.CurrentVolume) + float64(arg.MatchVolume)))
		investment.CapitalCost = int64(roundUpCapitalCost)

		// create transaction
		transaction, err := t.store.CreateTransaction(ctx, db.CreateTransactionParams{
			InvestmentID:    arg.InvestmentId,
			Ticker:          investment.Ticker,
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
		_, err = t.store.UpdateInvestmentWhenBuying(ctx, db.UpdateInvestmentWhenBuyingParams{
			ID:                   investment.ID,
			BuyTransactionVolume: transaction.MatchVolume,
			BuyTransactionValue:  transaction.MatchValue,
			CapitalCost:          investment.CapitalCost,
			TransactionFee:       transaction.Fee,
			TransactionTax:       transaction.Tax,
			Status:               db.InvestmentStatusActive,
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

func (t *transactionService) createSellingTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error) {
	result, err := t.store.ExecTx(ctx, func(query *db.Queries) (interface{}, error) {
		// get investment
		investment, err := query.GetInvestmentById(ctx, arg.InvestmentId)
		if err != nil {
			return db.Transaction{}, err
		}
		if investment.CurrentVolume < arg.MatchVolume {
			err = fmt.Errorf("current volume is less than match volume")
			return db.Transaction{}, err
		}

		// create a transaction
		transaction, err := t.store.CreateTransaction(ctx, db.CreateTransactionParams{
			InvestmentID:    arg.InvestmentId,
			Ticker:          investment.Ticker,
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
			Return:          (arg.MatchPrice-investment.CapitalCost)*arg.MatchVolume - arg.Fee - arg.Tax,
			Status:          arg.Status,
		})
		if err != nil {
			return db.Transaction{}, err
		}

		// create account's entry
		entry, err := query.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: investment.AccountID,
			Amount:    transaction.MatchValue - transaction.Fee - transaction.Tax,
			Type:      db.EntryTypeIT,
		})
		if err != nil {
			return db.Transaction{}, err
		}

		// update account's balance
		query.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			Amount: entry.Amount,
			ID:     investment.AccountID,
		})
		if investment.CurrentVolume-transaction.MatchVolume == 0 {
			investment.Status = db.InvestmentStatusSellout
		}

		// update investment when selling
		_, err = query.UpdateInvestmentWhenSeling(ctx, db.UpdateInvestmentWhenSelingParams{
			ID:                    transaction.InvestmentID,
			SellTransactionVolume: transaction.MatchVolume,
			SellTransactionValue:  transaction.MatchValue,
			TransactionFee:        transaction.Fee,
			TransactionTax:        transaction.Tax,
			UpdatedDate: pgtype.Timestamp{
				Time:  time.Now(),
				Valid: true,
			},
			Status: investment.Status,
		})
		if err != nil {
			return db.Transaction{}, err
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

func (t *transactionService) checkBuyingTransaction(ctx context.Context, investment *db.Investment, transaction *db.Transaction, account db.Account, q *db.Queries) (returnError int64, err error) {
	totalTransactionValue := (transaction.MatchPrice * transaction.MatchVolume) + transaction.Fee + transaction.Tax
	if account.Balance < totalTransactionValue {
		return 0, fmt.Errorf("account balance is less than transation cost")
	}

	// create entry
	entry, err := q.CreateEntry(ctx, db.CreateEntryParams{
		AccountID: investment.AccountID,
		Amount:    -totalTransactionValue,
		Type:      db.EntryTypeIT,
	})
	if err != nil {
		return 0, err
	}

	// add account balance
	account, err = q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
		ID:     account.ID,
		Amount: entry.Amount,
	})
	if err != nil {
		return 0, err
	}
	// calculate capital cost for each shares
	currentInvestmentValue := investment.CurrentVolume * investment.CapitalCost
	roundUpCapitalCost := math.Round(
		(float64(currentInvestmentValue) + float64(totalTransactionValue)) /
			(float64(investment.CurrentVolume) + float64(transaction.MatchVolume)))

	investment.CapitalCost = int64(roundUpCapitalCost)
	returnError = (investment.CapitalCost - transaction.Cost)
	absReturnError := int64(math.Abs(float64(returnError)))
	if absReturnError > 10 {
		return 0, fmt.Errorf("transaction cost error: %s - %s -- transaction cost: %v  -- investment cost: %v", transaction.Ticker, transaction.Trade, transaction.Cost, investment.CapitalCost)
	}

	// update investment
	updatedInvestment, err := q.UpdateInvestmentWhenBuying(ctx, db.UpdateInvestmentWhenBuyingParams{
		ID:                   investment.ID,
		BuyTransactionVolume: transaction.MatchVolume,
		BuyTransactionValue:  transaction.MatchValue,
		CapitalCost:          transaction.Cost,
		TransactionFee:       transaction.Fee,
		TransactionTax:       transaction.Tax,
		Status:               db.InvestmentStatusActive,
		UpdatedDate: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return 0, err
	}
	if updatedInvestment.BuyVolume-investment.BuyVolume != transaction.MatchVolume ||
		updatedInvestment.BuyValue-investment.BuyValue != transaction.MatchValue ||
		updatedInvestment.CapitalCost != transaction.Cost ||
		updatedInvestment.Fee-investment.Fee != transaction.Fee ||
		updatedInvestment.Tax-investment.Tax != transaction.Tax {
		return 0, fmt.Errorf("updated buy %s investment have incorrect datas", updatedInvestment.Ticker)
	}
	return returnError, nil
}

func (t *transactionService) checkSellingTransaction(ctx context.Context, investment *db.Investment, transaction *db.Transaction, q *db.Queries) error {
	if investment.CurrentVolume < transaction.MatchVolume {
		return fmt.Errorf("%s - %s transaction's match volume is lesser than investment volume", transaction.Ticker, transaction.Trade)
	}
	// check return
	returnValue := (transaction.MatchPrice-transaction.Cost)*transaction.MatchVolume - transaction.Fee - transaction.Tax
	if int64(math.Abs(float64(returnValue-transaction.Return))) > 0 {
		return fmt.Errorf("%s - %s transaction's return value is not match with excel data", transaction.Ticker, transaction.Trade)
	}

	// create account's entry
	entry, err := q.CreateEntry(ctx, db.CreateEntryParams{
		AccountID: investment.AccountID,
		Amount:    transaction.MatchValue - transaction.Fee - transaction.Tax,
		Type:      db.EntryTypeIT,
	})
	if err != nil {
		return err
	}

	// update account's balance
	q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
		Amount: entry.Amount,
		ID:     investment.AccountID,
	})

	// update investment
	if investment.CurrentVolume-transaction.MatchVolume == 0 {
		investment.Status = db.InvestmentStatusSellout
	}
	updatedInvestment, err := q.UpdateInvestmentWhenSeling(ctx, db.UpdateInvestmentWhenSelingParams{
		ID:                    investment.ID,
		SellTransactionVolume: transaction.MatchVolume,
		SellTransactionValue:  transaction.MatchValue,
		TransactionFee:        transaction.Fee,
		TransactionTax:        transaction.Tax,
		CapitalCost:           transaction.Cost,
		UpdatedDate: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		Status: investment.Status,
	})
	if err != nil {
		return err
	}
	if updatedInvestment.SellVolume-investment.SellVolume != transaction.MatchVolume ||
		updatedInvestment.SellValue-investment.SellValue != transaction.MatchValue ||
		updatedInvestment.Fee-investment.Fee != transaction.Fee ||
		updatedInvestment.Tax-investment.Tax != transaction.Tax ||
		updatedInvestment.UpdatedDate.Time.Truncate(24*time.Hour) != time.Now().UTC().Truncate(24*time.Hour) {
		return fmt.Errorf("updated sell %s investment have incorrect datas", updatedInvestment.Ticker)
	}
	return nil
}
