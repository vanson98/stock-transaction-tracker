package controller

import (
	"net/http"
	"strconv"
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
	var requestModel investment_model.SearchInvestmentRequestModel
	err := c.ShouldBindQuery(&requestModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var searchPram = db.SearchInvestmentPagingParams{
		AccountIds: requestModel.AccountIds,
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
		AccountIds: requestModel.AccountIds,
		SearchText: searchPram.SearchText,
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

func (ic *InvestmentController) GetById(c *gin.Context) {
	idParam, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, "id is required")
		return
	}
	investmentId, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	investment, err := ic.investmentService.GetById(c, int64(investmentId))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, investment)
}
