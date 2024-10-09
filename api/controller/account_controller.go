package controller

import (
	"net/http"
	apimodels "stt/api/models"
	db "stt/database/postgres/sqlc"
	"stt/domain"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	AccountService domain.IAccountService
}

func (ac *AccountController) CreateNewAccount(ctx *gin.Context) {
	var reqBody apimodels.CreateAccountRequest
	if err := ctx.ShouldBindBodyWithJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	account, err := ac.AccountService.CreateNew(ctx, db.CreateAccountParams{
		ChannelName: reqBody.ChannelName,
		Owner:       reqBody.Owner,
		Currency:    reqBody.Currency,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, account)
}
