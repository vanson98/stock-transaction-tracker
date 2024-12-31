package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
)

type IInvestmentService interface {
	Create(c context.Context, param db.CreateInvestmentParams) (db.Investment, error)
	GetById(c context.Context, id int64) (db.Investment, error)
	SearchPaging(c context.Context, param db.SearchInvestmentPagingParams) ([]db.InvestmentOverview, error)
	Count(c context.Context, db db.CountInvestmentParams) (int64, error)
}
