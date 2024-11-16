package controller

import (
	"fmt"
	"net/http"
	investment_model "stt/api/models/investment"
	db "stt/database/postgres/sqlc"
	sv_interface "stt/services/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type InvestmentController struct {
	investmentService sv_interface.IInvestmentService
}

func InitInvestmentController(investmentService sv_interface.IInvestmentService) InvestmentController {
	return InvestmentController{
		investmentService: investmentService,
	}
}

func (ic *InvestmentController) Search(c *gin.Context) {
	var requestModel investment_model.SearchInvestmentModel
	err := c.ShouldBindQuery(&requestModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var searchPram = db.SearchInvestmentPagingParams{
		AccountID:  requestModel.AccountId,
		SearchText: "%" + requestModel.SearchText + "%",
		OrderBy:    requestModel.OrderBy,
		SortType:   requestModel.SortType,
		TakeLimit:  int32(requestModel.PageSize),
		FromOffset: (requestModel.Page - 1) * requestModel.PageSize,
	}
	searchResult, err := ic.investmentService.SearchPaging(c, searchPram)
	if err != nil && err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	totalResult, err := ic.investmentService.Count(c, db.CountInvestmentParams{
		AccountID:  requestModel.AccountId,
		SearchText: requestModel.SearchText,
	})
	if err != nil && err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, investment_model.SearchInvestmentResponseModel{
		Investments: searchResult,
		TotalItems:  totalResult,
	})
}

func (ic *InvestmentController) Create(c *gin.Context) {
	var createInvestmentModel investment_model.CreateInvestmentModel
	err := c.ShouldBindBodyWithJSON(&createInvestmentModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// check investment exist
	ivm, err := ic.investmentService.GetByTicker(c, createInvestmentModel.Ticker)
	if err != nil && err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	} else if ivm.ID > 0 && err == nil {
		err := fmt.Errorf("investment already exist")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	investment, err := ic.investmentService.Create(c, db.CreateInvestmentParams{
		AccountID:   createInvestmentModel.AccountID,
		Ticker:      createInvestmentModel.Ticker,
		CompanyName: pgtype.Text{String: createInvestmentModel.CompanyName, Valid: true},
		Description: pgtype.Text{String: createInvestmentModel.Description, Valid: true},
		MarketPrice: createInvestmentModel.MarketPrice,
		Status:      db.InvestmentStatusInactive,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, investment)
}
