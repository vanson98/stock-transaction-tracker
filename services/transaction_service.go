package services

import (
	"context"
	db "stt/database/postgres/sqlc"
	sv_interface "stt/services/interfaces"
)

type transactionService struct {
	store db.IStore
}

func InitTransactionService(store db.IStore) sv_interface.ITransactionService {
	return transactionService{
		store: store,
	}
}

// CreateNew implements sv_interface.ITransactionService.
func (t transactionService) CreateNew(ctx context.Context, arg db.CreateTransactionParams) (db.Transaction, error) {
	return t.store.CreateTransaction(ctx, arg)
}
