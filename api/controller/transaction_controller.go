package controller

import (
	"net/http"
	transaction_model "stt/api/models/transaction"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
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

func (tc transactionController) GetPaging(c *gin.Context) {
	var requestModel transaction_model.GetTransactionsPagingModel
	err := c.ShouldBindQuery(&requestModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	transactions, err := tc.transactionService.GetPaging(c, db.GetTransactionsPagingParams{
		AccountIds: requestModel.AccountIds,
		ToLimit:    requestModel.PageSize,
		Ticker:     requestModel.Ticker,
		FromOffset: (requestModel.Page - 1) * requestModel.PageSize,
		OrderBy:    requestModel.OrderBy,
		OrderType:  requestModel.OrderType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	totalRow, err := tc.transactionService.CountTransaction(c, db.CountTransactionsParams{
		AccountIds: requestModel.AccountIds,
		Ticker:     requestModel.Ticker,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	result := transaction_model.GetTransactionsPagingResponseModel{
		Transactions: transactions,
		Total:        int32(totalRow),
	}
	c.JSON(http.StatusOK, result)
}

func (tc transactionController) CreateNewTransaction(c *gin.Context) {
	var requestModel transaction_model.CreateTransactionModel
	err := c.ShouldBindBodyWithJSON(&requestModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	tradingDate, err := time.Parse("2006-01-02T15:04:05Z", requestModel.TradingDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := tc.transactionService.AddTransaction(c, dtos.CreateTransactionDto{
		InvestmentId: requestModel.InvestmentID,
		TradingDate: pgtype.Timestamp{
			Time:  tradingDate,
			Valid: true,
		},
		Trade:       db.TradeType(requestModel.Trade),
		Volume:      requestModel.MatchVolume,
		OrderPrice:  requestModel.MatchPrice,
		MatchVolume: requestModel.MatchVolume,
		MatchPrice:  requestModel.MatchPrice,
		Fee:         requestModel.Fee,
		Tax:         requestModel.Tax,
		Status:      requestModel.Status,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, result)
}
