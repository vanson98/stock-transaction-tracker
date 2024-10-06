package domain

import (
	"context"
	db "stt/database/postgres/sqlc"
)

type IAccountService interface {
	CreateNew(ctx context.Context, param db.CreateAccountParams) (db.Account, error)
}
