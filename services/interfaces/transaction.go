package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
)

type ITransactionService interface {
	CreateNew(ctx context.Context, arg db.CreateTransactionParams) (db.Transaction, error)
}
