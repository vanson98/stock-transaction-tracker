package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
)

type IAccountService interface {
	CreateNew(ctx context.Context, param db.CreateAccountParams) (db.Account, error)
	GetAllByOwner(ctx context.Context, owner string) ([]db.GetAccountPagingRow, error)
	ListAllByOwner(ctx context.Context, owner string) ([]db.ListAllAccountRow, error)
	GetById(ctx context.Context, id int64) (db.Account, error)
	UpdateBalance(ctx context.Context, param db.AddAccountBalanceParams) (db.Account, error)
	TransferMoney(ctx context.Context, arg dtos.TransferMoneyTxParam) (dtos.TransferMoneyTxResult, error)
	Delete(ctx context.Context, accountId int64) error
	GetAccountInfoByIds(ctx context.Context, ids []int64) ([]db.GetAccountInfoByIdsRow, error)
}
