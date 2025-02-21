package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	transaction_model "stt/api/models/transaction"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
	sv_interface "stt/services/interfaces"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/xuri/excelize/v2"
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
	sumInfo, err := tc.transactionService.GetSummarizeInfo(c, db.GetTransactionSummarizeInfoParams{
		AccountIds: requestModel.AccountIds,
		Ticker:     requestModel.Ticker,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	result := transaction_model.GetTransactionsPagingResponseModel{
		Transactions:  transactions,
		TotalRow:      int32(sumInfo.TotalRows),
		SumMatchValue: sumInfo.SumMatchValue,
		SumFee:        sumInfo.SumFee,
		SumTax:        sumInfo.SumTax,
		SumReturn:     sumInfo.SumReturn,
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
	result, err := tc.transactionService.CreateTransaction(c, dtos.CreateTransactionDto{
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

func (ac *transactionController) ImportTransactions(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fileHeader, ok := form.File["tcbs_transaction_export_data"]
	account_id, ok := form.Value["account_id"]

	if !ok {
		c.JSON(http.StatusBadRequest, fmt.Errorf("bad request"))
		return
	}
	accountId, _ := strconv.Atoi(account_id[0])

	contentType := fileHeader[0].Header.Get("Content-Type")
	if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("content type is not in xlsx format"))
		return
	}

	exportFile, _ := fileHeader[0].Open()

	f, err := excelize.OpenReader(exportFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	defer func() {
		f.Close()
	}()
	rows, err := f.GetRows("Sheet 1")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	transactions := []db.Transaction{}
	for i, row := range rows {
		if i > 14 {
			if row == nil {
				break
			}
			tradingDate, err := time.Parse("02/01/2006", row[1])
			if err != nil {
				continue
			}
			var trade db.TradeType
			if row[2] == "Mua" {
				trade = db.TradeTypeBUY
			} else {
				trade = db.TradeTypeSELL
			}
			volume, _ := strconv.Atoi(strings.Replace(row[3], ",", "", -1))
			orderPrice, _ := strconv.Atoi(strings.Replace(row[4], ",", "", -1))
			matchVolume, _ := strconv.Atoi(strings.Replace(row[5], ",", "", -1))
			matchPrice, _ := strconv.Atoi(strings.Replace(row[6], ",", "", -1))
			matchValue, _ := strconv.Atoi(strings.Replace(row[7], ",", "", -1))
			fee, _ := strconv.Atoi(strings.Replace(row[8], ",", "", -1))
			tax, _ := strconv.Atoi(strings.Replace(row[9], ",", "", -1))
			cost, _ := strconv.Atoi(strings.Replace(row[10], ",", "", -1))
			returnValue, _ := strconv.Atoi(strings.Replace(row[11], ",", "", -1))
			transaction := db.Transaction{
				Ticker: row[0],
				TradingDate: pgtype.Timestamp{
					Time:  tradingDate,
					Valid: true,
				},
				Trade:       trade,
				Volume:      int64(volume),
				OrderPrice:  int64(orderPrice),
				MatchVolume: int64(matchVolume),
				MatchPrice:  int64(matchPrice),
				MatchValue:  int64(matchValue),
				Fee:         int64(fee),
				Tax:         int64(tax),
				Cost:        int64(cost),
				Return:      int64(returnValue),
				Status:      db.TransactionStatusCOMPLETED,
			}
			transactions = append(transactions, transaction)
		}
	}
	result, err := ac.transactionService.ImportTransactions(c, int64(accountId), transactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ac *transactionController) AddTransaction(c *gin.Context) {
	var addTransactionRequest transaction_model.AddTransactionRequest
	err := c.ShouldBindBodyWithJSON(&addTransactionRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	tradingDate, err := time.Parse("02/01/2006", addTransactionRequest.TradingDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	transaction, err := ac.transactionService.ImportTransaction(c, addTransactionRequest.AccountId, db.Transaction{
		Ticker: addTransactionRequest.Ticker,
		TradingDate: pgtype.Timestamp{
			Time:  tradingDate,
			Valid: true,
		},
		Trade:       db.TradeType(addTransactionRequest.Trade),
		Volume:      addTransactionRequest.Volume,
		OrderPrice:  addTransactionRequest.OrderPrice,
		MatchVolume: addTransactionRequest.MatchVolume,
		MatchPrice:  addTransactionRequest.MatchPrice,
		MatchValue:  addTransactionRequest.MatchValue,
		Fee:         addTransactionRequest.Fee,
		Tax:         addTransactionRequest.Tax,
		Cost:        addTransactionRequest.Cost,
		Return:      addTransactionRequest.Return,
		Status:      db.TransactionStatus(addTransactionRequest.Status),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
	} else {
		// response to cron job
		c.JSON(http.StatusOK, transaction)
	}

	apiGatewayNoti := struct {
		TransactionId int64
		Error         string
	}{
		Error:         "",
		TransactionId: transaction.ID,
	}
	if err != nil {
		apiGatewayNoti.Error = err.Error()
	}
	// transer result to api gate way
	client := http.Client{}
	jsonData, _ := json.Marshal(apiGatewayNoti)
	buffer := bytes.Buffer{}
	buffer.Write(jsonData)
	_, err = client.Post("http://localhost:6061/notification", "application/json", &buffer)
	if err != nil {
		fmt.Println(err)
	}
}
