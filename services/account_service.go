package services

import (
	"context"
	"fmt"
	db "stt/database/postgres/sqlc"
	"stt/domain"
	"stt/services/dtos"
	"time"
)

type accountService struct {
	store   db.IStore
	timeout time.Duration
}

func InitAccountService(store db.IStore, timeout time.Duration) domain.IAccountService {
	return accountService{
		store:   store,
		timeout: timeout,
	}
}

// CreateNew implements domain.IAccountService.
func (as accountService) CreateNew(ctx context.Context, param db.CreateAccountParams) (db.Account, error) {
	return as.store.CreateAccount(ctx, param)
}

// GetAllPaging implements domain.IAccountService.
func (as accountService) GetAllPaging(ctx context.Context, param db.ListAccountsParams) ([]db.Account, error) {
	return as.store.ListAccounts(ctx, param)
}

// UpdateBalance implements domain.IAccountService.
func (as accountService) UpdateBalance(ctx context.Context, param db.AddAccountBalanceParams) (db.Account, error) {
	return as.store.AddAccountBalance(ctx, param)
}

var txKey = struct{}{}

// TransferMoney implements domain.IAccountService.
func (as accountService) TransferMoney(ctx context.Context, arg dtos.TransferMoneyTxParam) (dtos.TransferMoneyTxResult, error) {
	var result dtos.TransferMoneyTxResult

	err := as.store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		txName := ctx.Value(txKey)

		// create a entry
		fmt.Println(txName, "create a entry")
		accEntry, err := q.CreateEntry(ctx, db.CreateEntryParams{
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
		result.UpdatedAccount, err = q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
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

// Delete implements domain.IAccountService.
func (as accountService) Delete(ctx context.Context, accountId int64) error {
	return as.store.DeleteAccount(ctx, accountId)
}
