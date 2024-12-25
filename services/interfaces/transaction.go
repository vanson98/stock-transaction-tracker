package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
)

type ITransactionService interface {
	CreateTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error)
	GetById(ctx context.Context, id int64) (db.Transaction, error)
	GetPaging(ctx context.Context, param db.GetTransactionsPagingParams) ([]db.GetTransactionsPagingRow, error)
	CountTransaction(ctx context.Context, param db.CountTransactionsParams) (int64, error)
	InsertTransaction(ctx context.Context, accountId int64, transactions []db.Transaction) (bool, error)
}
