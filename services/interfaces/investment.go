package sv_interface

import (
	"context"
	db "stt/database/postgres/sqlc"
)

type IInvestmentService interface {
	Create(c context.Context, param db.CreateInvestmentParams) (db.Investment, error)
	//Update(c context.Context, investment *Investment)
	GetById(c context.Context, id int32)
	GetAll(c context.Context)
	Delete(c context.Context, id int32)
	GetByTicker(ctx context.Context, ticker string) (db.Investment, error)
}
