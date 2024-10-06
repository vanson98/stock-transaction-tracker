package services

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/domain"
	"time"
)

type accountService struct {
	store   *db.Store
	timeout time.Duration
}

func InitAccountService(store *db.Store, timeout time.Duration) domain.IAccountService {
	return accountService{
		store:   store,
		timeout: timeout,
	}
}

// CreateNew implements domain.IAccountService.
func (as accountService) CreateNew(ctx context.Context, param db.CreateAccountParams) (db.Account, error) {
	return as.store.Queries.CreateAccount(ctx, param)
}
