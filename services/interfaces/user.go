package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
)

type IUserService interface {
	CreateNew(ctx context.Context, param db.CreateUserParams) (db.User, error)
	GetByUserName(ctx context.Context, username string) (db.User, error)
}
