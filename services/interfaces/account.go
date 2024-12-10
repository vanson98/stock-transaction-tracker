package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
)

type IAccountService interface {
	CreateNew(ctx context.Context, param db.CreateAccountParams) (db.Account, error)
	GetAllByOwner(ctx context.Context, owner string) ([]db.GetAccountPagingRow, error)
	ListAllAccount(ctx context.Context) ([]db.ListAllAccountRow, error)
	GetById(ctx context.Context, id int64) (db.Account, error)
	UpdateBalance(ctx context.Context, param db.AddAccountBalanceParams) (db.Account, error)
	TransferMoney(ctx context.Context, arg dtos.TransferMoneyTxParam) (dtos.TransferMoneyTxResult, error)
	Delete(ctx context.Context, accountId int64) error
	GetAccountInfoById(ctx context.Context, id int64) (db.GetAccountInfoByIdRow, error)
}
