package services

import (
	"context"
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

func (as accountService) GetById(ctx context.Context, id int64) (db.Account, error) {
	return as.store.GetAccountById(ctx, id)
}

func (as accountService) GetStockAssetOverview(ctx context.Context, ids []int64) ([]db.GetStockAssetOverviewRow, error) {
	return as.store.GetStockAssetOverview(ctx, ids)
}

// ListAllAccount implements sv_interface.IAccountService.
func (as accountService) ListAllByOwner(ctx context.Context, owner string) ([]db.ListAllAccountRow, error) {
	return as.store.ListAllAccount(ctx, owner)
}

// GetAllPaging implements sv_interface.IAccountService.
func (as accountService) GetAllOverview(ctx context.Context, owner string) ([]db.GetAllAccountOverviewRow, error) {
	return as.store.GetAllAccountOverview(ctx, owner)
}

var TxKey = struct{}{}

// TransferMoney implements sv_interface.IAccountService.
func (as accountService) TransferMoney(ctx context.Context, arg dtos.TransferMoneyTxParam) (dtos.TransferMoneyTxResult, error) {
	var result dtos.TransferMoneyTxResult

	_, err := as.store.ExecTx(ctx, func(q *db.Queries) (interface{}, error) {
		var err error
		//txName := ctx.Value(TxKey)

		// create a entry
		//fmt.Println(txName, "create a entry")
		accEntry, err := q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: arg.AccountID,
			Type:      arg.EntryType,
			Amount:    arg.Amount,
		})
		if err != nil {
			return nil, err
		}
		result.Entry = accEntry

		//update account balance
		//fmt.Println(txName, "update account balance")
		result.UpdatedAccount, err = q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			ID:     arg.AccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return result, err
}
