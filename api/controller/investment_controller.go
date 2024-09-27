package controller

import (
	db "stt/database/postgres/sqlc"
	"stt/domain"

	"github.com/gin-gonic/gin"
)

type InvestmentController struct {
	InvestmentService domain.IInvestmentService
}

func (ic *InvestmentController) GetAll(c *gin.Context) {
	ic.InvestmentService.GetAll(c)
}

func (ic *InvestmentController) Create(c *gin.Context) {
	newInvestment := db.CreateInvestmentParams{}
	ic.InvestmentService.Create(c, newInvestment)
}
