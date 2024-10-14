package services

import (
	"context"
	db "stt/database/postgres/sqlc"
	sv_interface "stt/services/interfaces"
	"time"
)

type investmentService struct {
	store          db.IStore
	contextTimeout time.Duration
}

func InitInvestmentService(store db.IStore, timeout time.Duration) sv_interface.IInvestmentService {
	return investmentService{
		store:          store,
		contextTimeout: timeout,
	}
}

// Create implements domain.IInvestmentService.
func (ivs investmentService) Create(c context.Context, param db.CreateInvestmentParams) (db.Investment, error) {
	ctx, cancel := context.WithTimeout(c, ivs.contextTimeout)
	defer cancel()
	return ivs.store.CreateInvestment(ctx, param)
}

// Delete implements domain.IInvestmentService.
func (i investmentService) Delete(c context.Context, id int32) {
	panic("unimplemented")
}

// GetAll implements domain.IInvestmentService.
func (i investmentService) GetAll(c context.Context) {
	panic("unimplemented")
}

// GetById implements domain.IInvestmentService.
func (i investmentService) GetById(c context.Context, id int32) {
	panic("unimplemented")
}

// Update implements domain.IInvestmentService.
func (i investmentService) Update(c context.Context, investment *db.Investment) {
	panic("unimplemented")
}
