package controller

import (
	"net/http"
	transaction_model "stt/api/models/transaction"
	db "stt/database/postgres/sqlc"
	sv_interface "stt/services/interfaces"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type transactionController struct {
	transactionService sv_interface.ITransactionService
}

func InitTransactionController(transService sv_interface.ITransactionService) transactionController {
	return transactionController{
		transactionService: transService,
	}
}

func (tc transactionController) CreateNewTransaction(c *gin.Context) {
	var requestModel transaction_model.CreateTransactionModel
	err := c.ShouldBindBodyWithJSON(&requestModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	tradingDate, err := time.Parse("2006-01-02 15:04:05", requestModel.TradingDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := tc.transactionService.CreateNew(c, db.CreateTransactionParams{
		InvestmentID: requestModel.InvestmentID,
		Ticker:       requestModel.Ticker,
		TradingDate: pgtype.Timestamp{
			Time:  tradingDate,
			Valid: true,
		},
		Trade:           requestModel.Trade,
		Volume:          requestModel.Volume,
		OrderPrice:      requestModel.OrderPrice,
		MatchVolume:     requestModel.MatchVolume,
		MatchPrice:      requestModel.MatchPrice,
		MatchValue:      requestModel.MatchValue,
		Fee:             requestModel.Fee,
		Tax:             requestModel.Tax,
		Cost:            requestModel.Cost,
		CostOfGoodsSold: requestModel.CostOfGoodsSold,
		Return:          requestModel.Return,
		Status:          requestModel.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, result)
}
