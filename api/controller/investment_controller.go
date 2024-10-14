package controller

import (
	sv_interface "stt/services/interfaces"
)

type InvestmentController struct {
	InvestmentService sv_interface.IInvestmentService
}

// func (ic *InvestmentController) GetAll(c *gin.Context) {
// 	ic.InvestmentService.GetAll(c)
// }

// func (ic *InvestmentController) Create(c *gin.Context) {
// 	newInvestment := db.CreateInvestmentParams{}
// 	ic.InvestmentService.Create(c, newInvestment)
// }
