package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	account_model "stt/api/models/account"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
	sv_interface "stt/services/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AccountController struct {
	AccountService sv_interface.IAccountService
}

func (ac *AccountController) CreateNewAccount(ctx *gin.Context) {
	var reqBody account_model.CreateAccountRequest
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
	requestParam := account_model.GetAccountRequest{}
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

func (ac *AccountController) GetAccountInfoByIds(ctx *gin.Context) {
	var requestData account_model.GetAccountInfoRequest
	err := ctx.ShouldBindQuery(&requestData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	acc, err := ac.AccountService.GetAccountInfoByIds(ctx, requestData.Ids)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, acc)
}

func (ac *AccountController) GetListAccount(ctx *gin.Context) {
	owner := ctx.Request.URL.Query().Get("owner")
	if owner == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("required owner")))
		return
	}
	accounts, err := ac.AccountService.ListAllByOwner(ctx, owner)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

func (ac *AccountController) GetAccoutPaging(ctx *gin.Context) {
	requestDataModel := account_model.SearchAccountRequest{}
	err := ctx.ShouldBindQuery(&requestDataModel)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accounts, err := ac.AccountService.GetAllByOwner(ctx, requestDataModel.Onwer)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	ctx.JSON(http.StatusOK, accounts)
}

func (ac *AccountController) TransferMoney(ctx *gin.Context) {
	transferRequest := account_model.TransferMoneyRequest{}
	if err := ctx.ShouldBindBodyWithJSON(&transferRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !ac.validAccount(ctx, transferRequest.AccountID, transferRequest.Currency) {
		return
	}

	if !ac.validTransfer(ctx, transferRequest.AccountID, transferRequest.Amount) {
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

func (ac *AccountController) validTransfer(ctx *gin.Context, accountId int64, amount int64) bool {
	updateAccount, err := ac.AccountService.GetById(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if updateAccount.Balance+amount < 0 {
		err := fmt.Errorf("can not perform this tranfer because balance is negative")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
