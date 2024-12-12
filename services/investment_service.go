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

// CountSearchPaging implements sv_interface.IInvestmentService.
func (ivs investmentService) Count(c context.Context, param db.CountInvestmentParams) (int64, error) {
	return ivs.store.CountInvestment(c, param)
}

// SearchPaging implements sv_interface.IInvestmentService.
func (ivs investmentService) SearchPaging(c context.Context, param db.SearchInvestmentPagingParams) ([]db.SearchInvestmentPagingRow, error) {
	return ivs.store.SearchInvestmentPaging(c, param)
}

// Create implements domain.IInvestmentService.
func (ivs investmentService) Create(c context.Context, param db.CreateInvestmentParams) (db.Investment, error) {
	ctx, cancel := context.WithTimeout(c, ivs.contextTimeout)
	defer cancel()
	return ivs.store.CreateInvestment(ctx, param)
}

// Delete implements domain.IInvestmentService.
func (i investmentService) Delete(c context.Context, id int64) {
	panic("unimplemented")
}

// GetByTicker implements sv_interface.IInvestmentService.
func (ivs investmentService) GetByTicker(ctx context.Context, ticker string) (db.Investment, error) {
	return ivs.store.GetInvestmentByTicker(ctx, ticker)
}

// GetAll implements domain.IInvestmentService.
func (i investmentService) GetAll(c context.Context) {
	panic("unimplemented")
}

// GetById implements domain.IInvestmentService.
func (i investmentService) GetById(c context.Context, id int64) (db.Investment, error) {
	return i.store.GetInvestmentById(c, id)
}

// Update implements domain.IInvestmentService.
func (i investmentService) Update(c context.Context, investment *db.Investment) {
	panic("unimplemented")
}
