package controller

import (
	"net/http"
	db "stt/database/postgres/sqlc"
	"stt/domain"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	AccountService domain.IAccountService
}

type createAccountRequest struct {
	ChannelName string `json:"channel_name" binding:"required"`
	Owner       string `json:"owner" binding:"required"`
	Currency    string `json:"currency" binding:"required,oneof=USD VND EUR"`
}

func (ac *AccountController) CreateNewAccount(ctx *gin.Context) {
	var reqBody createAccountRequest
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
