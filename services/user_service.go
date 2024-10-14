package services

import (
	"context"
	db "stt/database/postgres/sqlc"
	sv_interface "stt/services/interfaces"
)

type userService struct {
	store db.IStore
}

func InitUserService(store db.IStore) sv_interface.IUserService {
	return userService{
		store: store,
	}
}

// CreateNew implements sv_interface.IUserService.
func (us userService) CreateNew(ctx context.Context, param db.CreateUserParams) (db.User, error) {
	return us.store.CreateUser(ctx, param)
}

// GetByUserName implements sv_interface.IUserService.
func (us userService) GetByUserName(ctx context.Context, username string) (db.User, error) {
	return us.store.GetUser(ctx, username)
}
