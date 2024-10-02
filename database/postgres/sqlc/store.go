package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	dbPool  *pgxpool.Pool
	queries *Queries
}

func NewStore(dbp *pgxpool.Pool) *Store {
	return &Store{
		dbPool:  dbp,
		queries: New(dbp),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.dbPool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

type TransferTxParam struct {
	AccountID int64     `json:"account"`
	Amount    int64     `json:"amount"`
	EntryType EntryType `json:"entry_from_type"`
}

type TransferTxResult struct {
	UpdatedAccount Account `json:"account"`
	AccountEntry   Entry   `json:"account_entry"`
}

var txKey = struct{}{}

// TransferTx perform a money transfer in or out of account
// Create account entries and update account's balance within a single database transaction
func (store *Store) TranserTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)

		// create a entry
		fmt.Println(txName, "create a entry")
		accEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.AccountID,
			Type:      arg.EntryType,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		result.AccountEntry = accEntry

		// get account for update
		// fmt.Println(txName, "get account for update")
		// account, err := q.GetAccountForUpdate(ctx, arg.AccountID)
		// if err != nil {
		// 	return err
		// }

		//update account balance
		fmt.Println(txName, "update account balance")
		result.UpdatedAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.AccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}