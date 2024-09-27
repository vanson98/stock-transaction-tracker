package domain

import (
	"context"
	db "stt/database/postgres/sqlc"
)

type IInvestmentRepository interface {
	Create(c context.Context, investmentParam db.CreateInvestmentParams) (db.Investment, error)
	//Update(c context.Context, investment *Investment)
	GetById(c context.Context, id int32)
	GetAll(c context.Context)
	Delete(c context.Context, id int32)
}

type IInvestmentService interface {
	Create(c context.Context, param db.CreateInvestmentParams) (db.Investment, error)
	//Update(c context.Context, investment *Investment)
	GetById(c context.Context, id int32)
	GetAll(c context.Context)
	Delete(c context.Context, id int32)
}
