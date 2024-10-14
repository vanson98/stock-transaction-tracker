package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	apimodels "stt/api/models"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
	sv_interface "stt/services/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type AccountController struct {
	AccountService sv_interface.IAccountService
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
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505", "23503":
				ctx.JSON(http.StatusForbidden, errorResponse(pgErr))
				return
			}
		}
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

	if !ac.validAccount(ctx, transferRequest.AccountID, transferRequest.Currency) {
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

func (ac *AccountController) validAccount(ctx *gin.Context, accountId int64, currency string) bool {
	acc, err := ac.AccountService.GetById(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if acc.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", acc.ID, acc.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
