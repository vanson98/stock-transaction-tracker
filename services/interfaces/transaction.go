package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
)

type ITransactionService interface {
	CreateTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error)
	InsertTransaction(ctx context.Context, accountId int64, transactions []db.Transaction) (bool, error)
	GetPaging(ctx context.Context, param db.GetTransactionsPagingParams) ([]db.GetTransactionsPagingRow, error)
	GetSummarizeInfo(ctx context.Context, param db.GetTransactionSummarizeInfoParams) (db.GetTransactionSummarizeInfoRow, error)
}
