package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
)

type ITransactionService interface {
	CreateTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error)
	ImportTransaction(ctx context.Context, accountId int64, transaction db.Transaction, checkCost bool) (db.Transaction, error)
	GetPaging(ctx context.Context, param db.GetTransactionsPagingParams) ([]db.GetTransactionsPagingRow, error)
	GetSummarizeInfo(ctx context.Context, param db.GetTransactionSummarizeInfoParams) (db.GetTransactionSummarizeInfoRow, error)
}
