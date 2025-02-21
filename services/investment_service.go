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

func (ivs investmentService) Create(c context.Context, param db.CreateInvestmentParams) (db.Investment, error) {
	ctx, cancel := context.WithTimeout(c, ivs.contextTimeout)
	defer cancel()
	return ivs.store.CreateInvestment(ctx, param)
}

// GetById implements sv_interface.IInvestmentService.
func (ivs investmentService) GetById(c context.Context, id int64) (db.Investment, error) {
	return ivs.store.GetInvestmentById(c, id)
}

func (i investmentService) GetOverviewById(c context.Context, id int64) (db.InvestmentOverview, error) {
	return i.store.GetInvestmentOverviewById(c, id)
}

// CountSearchPaging implements sv_interface.IInvestmentService.
func (ivs investmentService) Count(c context.Context, param db.CountInvestmentParams) (int64, error) {
	return ivs.store.CountInvestment(c, param)
}

func (ivs investmentService) SearchPaging(c context.Context, param db.SearchInvestmentPagingParams) ([]db.InvestmentOverview, error) {
	return ivs.store.SearchInvestmentPaging(c, param)
}

// UpdateMarketPrice implements sv_interface.IInvestmentService.
func (ivs investmentService) UpdateMarketPrice(c context.Context, params db.UpdateMarketPriceParams) (db.UpdateMarketPriceRow, error) {
	return ivs.store.UpdateMarketPrice(c, params)
}
