package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
)

type ITransactionService interface {
	AddTransaction(ctx context.Context, arg dtos.CreateTransactionDto) (db.Transaction, error)
	GetById(ctx context.Context, id int64) (db.Transaction, error)
}
