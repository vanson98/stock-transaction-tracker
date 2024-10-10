package controller

import (
	"net/http"
	apimodels "stt/api/models"
	db "stt/database/postgres/sqlc"
	"stt/domain"
	"stt/services/dtos"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	AccountService domain.IAccountService
}

func (ac *AccountController) CreateNewAccount(ctx *gin.Context) {
	var reqBody apimodels.CreateAccountRequest
	if err := ctx.ShouldBindBodyWithJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := ac.AccountService.CreateNew(ctx, db.CreateAccountParams{
		ChannelName: reqBody.ChannelName,
		Owner:       reqBody.Owner,
		Currency:    reqBody.Currency,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (ac *AccountController) GetAccountById(ctx *gin.Context) {
	requestParam := apimodels.GetAccountRequest{}
	err := ctx.ShouldBindUri(&requestParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := ac.AccountService.GetById(ctx, requestParam.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

func (ac *AccountController) TransferMoney(ctx *gin.Context) {
	transferRequest := apimodels.TransferMoneyRequest{}
	if err := ctx.ShouldBindBodyWithJSON(&transferRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := ac.AccountService.TransferMoney(ctx, dtos.TransferMoneyTxParam{
		AccountID: transferRequest.AccountID,
		Amount:    transferRequest.Amount,
		EntryType: db.EntryType(transferRequest.EntryType),
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
