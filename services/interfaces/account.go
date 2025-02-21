package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
)

type IAccountService interface {
	CreateNew(ctx context.Context, param db.CreateAccountParams) (db.Account, error)
	GetById(ctx context.Context, id int64) (db.Account, error)
	GetStockAssetOverview(ctx context.Context, ids []int64) ([]db.GetStockAssetOverviewRow, error)
	ListAllByOwner(ctx context.Context, owner string) ([]db.ListAllAccountRow, error)
	GetAllOverview(ctx context.Context, owner string) ([]db.GetAllAccountOverviewRow, error)
	TransferMoney(ctx context.Context, arg dtos.TransferMoneyTxParam) (dtos.TransferMoneyTxResult, error)
}
