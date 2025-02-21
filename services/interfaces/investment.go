package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
)

type IInvestmentService interface {
	Create(c context.Context, param db.CreateInvestmentParams) (db.Investment, error)
	GetOverviewById(c context.Context, id int64) (db.InvestmentOverview, error)
	GetById(c context.Context, id int64) (db.Investment, error)
	SearchPaging(c context.Context, param db.SearchInvestmentPagingParams) ([]db.InvestmentOverview, error)
	Count(c context.Context, db db.CountInvestmentParams) (int64, error)
	UpdateMarketPrice(c context.Context, db db.UpdateMarketPriceParams) (db.UpdateMarketPriceRow, error)
}
