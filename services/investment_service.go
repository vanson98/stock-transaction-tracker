package services

import (
	"context"
	db "stt/database/postgres/sqlc"
	"stt/domain"
	"time"
)

type investmentService struct {
	queries        *db.Queries
	contextTimeout time.Duration
}

func InitInvestmentService(q *db.Queries, timeout time.Duration) domain.IInvestmentService {
	return investmentService{
		queries:        q,
		contextTimeout: timeout,
	}
}

// Create implements domain.IInvestmentService.
func (ivs investmentService) Create(c context.Context, param db.CreateInvestmentParams) (db.Investment, error) {
	ctx, cancel := context.WithTimeout(c, ivs.contextTimeout)
	defer cancel()
	return ivs.queries.CreateInvestment(ctx, param)
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
