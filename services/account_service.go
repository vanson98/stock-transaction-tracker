package services

import (
	"context"
	"fmt"
	db "stt/database/postgres/sqlc"

	"stt/services/dtos"
	sv_interface "stt/services/interfaces"
	"time"
)

type accountService struct {
	store   db.IStore
	timeout time.Duration
}

func InitAccountService(store db.IStore, timeout time.Duration) sv_interface.IAccountService {
	return accountService{
		store:   store,
		timeout: timeout,
	}
}

// CreateNew implements sv_interface.IAccountService.
func (as accountService) CreateNew(ctx context.Context, param db.CreateAccountParams) (db.Account, error) {
	return as.store.CreateAccount(ctx, param)
}

// ListAllAccount implements sv_interface.IAccountService.
func (as accountService) ListAllAccount(ctx context.Context) ([]db.Account, error) {
	return as.store.ListAllAccount(ctx)
}

// GetById implements sv_interface.IAccountService.
func (as accountService) GetById(ctx context.Context, id int64) (db.Account, error) {
	return as.store.GetAccountById(ctx, id)
}

// GetAllPaging implements sv_interface.IAccountService.
func (as accountService) GetAllPaging(ctx context.Context, param db.GetAccountsPagingParams) ([]db.Account, error) {
	return as.store.GetAccountsPaging(ctx, param)
}

// UpdateBalance implements sv_interface.IAccountService.
func (as accountService) UpdateBalance(ctx context.Context, param db.AddAccountBalanceParams) (db.Account, error) {
	return as.store.AddAccountBalance(ctx, param)
}

var TxKey = struct{}{}

// TransferMoney implements sv_interface.IAccountService.
func (as accountService) TransferMoney(ctx context.Context, arg dtos.TransferMoneyTxParam) (dtos.TransferMoneyTxResult, error) {
	var result dtos.TransferMoneyTxResult

	err := as.store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		txName := ctx.Value(TxKey)

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
		result.Entry = accEntry

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

// Delete implements sv_interface.IAccountService.
func (as accountService) Delete(ctx context.Context, accountId int64) error {
	return as.store.DeleteAccount(ctx, accountId)
}
