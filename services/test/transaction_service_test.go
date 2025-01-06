package service_test

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	db "stt/database/postgres/sqlc"
	"stt/services/dtos"
	"stt/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

func createRandomTransaction(t *testing.T, investmentId int64, ticker string, trade db.TradeType, matchVolume int64) db.Transaction {
	// create transaction params
	param := dtos.CreateTransactionDto{
		InvestmentId: investmentId,
		TradingDate: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		Trade:       trade,
		Volume:      util.RandomInt(100, 5000),
		OrderPrice:  util.RandomInt(10000, 100000),
		MatchVolume: matchVolume,
		MatchPrice:  util.RandomInt(10000, 100000),
		Fee:         util.RandomInt(10000, 20000),
		Tax:         util.RandomInt(10000, 20000),
		Status:      db.TransactionStatusCOMPLETED,
	}

	transaction, err := tranService.CreateTransaction(context.Background(), param)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)
	require.Equal(t, transaction.InvestmentID, investmentId)
	require.Equal(t, transaction.Ticker, ticker)
	require.Equal(t, transaction.TradingDate.Time.Truncate(24*time.Hour), param.TradingDate.Time.UTC().Truncate(24*time.Hour))
	require.Equal(t, transaction.Trade, param.Trade)
	require.Equal(t, transaction.Volume, param.Volume)
	require.Equal(t, transaction.OrderPrice, param.OrderPrice)
	require.Equal(t, transaction.MatchVolume, param.MatchVolume)
	require.Equal(t, transaction.Fee, param.Fee)
	require.Equal(t, transaction.Tax, param.Tax)
	require.Equal(t, transaction.Status, param.Status)
	require.Equal(t, param.Trade, transaction.Trade)
	require.Equal(t, transaction.MatchValue, param.MatchPrice*param.MatchVolume)

	return transaction
}

func TestCreateCircleTransaction(t *testing.T) {
	user := createRandomUser(t)
	acc := createRandomAccount(t, user.Username)
	investment := createDefaultInvestmnet(t, acc.ID)
	buyVolume := util.RandomInt(500, 1000)
	createRandomBuyTransaction(t, acc.ID, investment.ID, buyVolume)
	sellVolume := util.RandomInt(100, 500)
	createRandomSellTransaction(t, acc.ID, investment.ID, sellVolume)
}

func createRandomBuyTransaction(t *testing.T, accountId int64, investmentId int64, matchVolume int64) {
	investment, err := investmentService.GetById(context.Background(), investmentId)
	require.NoError(t, err)
	account, err := accService.GetById(context.Background(), accountId)
	require.NoError(t, err)
	transaction := createRandomTransaction(t, investment.ID, investment.Ticker, db.TradeTypeBUY, matchVolume)

	// check account's balance
	updatedAccount, err := accService.GetById(context.Background(), account.ID)
	account.Balance -= (transaction.MatchValue + transaction.Fee + transaction.Tax)
	require.NoError(t, err)
	require.Equal(t, updatedAccount.Balance, account.Balance)

	// check investment
	dbInvestment, err := investmentService.GetById(context.Background(), investment.ID)
	require.NoError(t, err)
	require.Equal(t, investment.Ticker, dbInvestment.Ticker)
	require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), dbInvestment.UpdatedDate.Time.Truncate(24*time.Hour))
	require.Equal(t, transaction.Cost, dbInvestment.CapitalCost)
	require.Equal(t, transaction.CostOfGoodsSold, dbInvestment.CapitalCost*transaction.MatchVolume)

	investment.BuyValue += transaction.MatchValue
	investment.BuyVolume += transaction.MatchVolume
	investment.CurrentVolume += transaction.MatchVolume
	require.Equal(t, investment.BuyValue, dbInvestment.BuyValue)
	require.Equal(t, investment.BuyVolume, dbInvestment.BuyVolume)
	require.Equal(t, transaction.Return, int64(0))

	investment.Fee += transaction.Fee
	investment.Tax += transaction.Tax
	require.Equal(t, investment.CurrentVolume, dbInvestment.CurrentVolume)
	require.Equal(t, investment.Fee, dbInvestment.Fee)
	require.Equal(t, investment.Tax, dbInvestment.Tax)
	require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), dbInvestment.UpdatedDate.Time.Truncate(24*time.Hour))

	if dbInvestment.CurrentVolume > 0 {
		require.Equal(t, db.InvestmentStatusActive, dbInvestment.Status)
	}
}

func createRandomSellTransaction(t *testing.T, accountId int64, investmentId int64, matchVolume int64) {
	investment, err := investmentService.GetById(context.Background(), investmentId)
	require.NoError(t, err)
	account, err := accService.GetById(context.Background(), accountId)
	require.NoError(t, err)
	transaction := createRandomTransaction(t, investment.ID, investment.Ticker, db.TradeTypeSELL, matchVolume)

	// check account's balance
	updatedAccount, err := accService.GetById(context.Background(), account.ID)
	account.Balance += (transaction.MatchValue - transaction.Fee - transaction.Tax)
	require.NoError(t, err)
	require.Equal(t, updatedAccount.Balance, account.Balance)

	// check investment
	updatedInvestment, err := investmentService.GetById(context.Background(), investment.ID)
	require.NoError(t, err)
	require.Equal(t, investment.Ticker, updatedInvestment.Ticker)
	require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), updatedInvestment.UpdatedDate.Time.Truncate(24*time.Hour))
	require.Equal(t, transaction.Cost, updatedInvestment.CapitalCost)
	require.Equal(t, transaction.CostOfGoodsSold, updatedInvestment.CapitalCost*transaction.MatchVolume)

	investment.SellValue += transaction.MatchValue
	investment.SellVolume += transaction.MatchVolume
	investment.CurrentVolume -= transaction.MatchVolume
	require.Equal(t, investment.SellValue, updatedInvestment.SellValue)
	require.Equal(t, investment.SellVolume, updatedInvestment.SellVolume)
	returnValue := ((transaction.MatchPrice - updatedInvestment.CapitalCost) * transaction.MatchVolume) - transaction.Fee - transaction.Tax
	require.Equal(t, transaction.Return, returnValue)

	investment.Fee += transaction.Fee
	investment.Tax += transaction.Tax
	require.Equal(t, investment.CurrentVolume, updatedInvestment.CurrentVolume)
	require.Equal(t, investment.Fee, updatedInvestment.Fee)
	require.Equal(t, investment.Tax, updatedInvestment.Tax)
	require.Equal(t, time.Now().UTC().Truncate(24*time.Hour), updatedInvestment.UpdatedDate.Time.Truncate(24*time.Hour))

	if updatedInvestment.CurrentVolume == 0 {
		require.Equal(t, db.InvestmentStatusSellout, updatedInvestment.Status)
	}
}

func TestImportTransaction(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user.Username)
	account2 := createRandomAccount(t, user.Username)
	account3 := createRandomAccount(t, user.Username)
	_, err := accService.TransferMoney(context.Background(), dtos.TransferMoneyTxParam{
		AccountID: account2.ID,
		Amount:    -account2.Balance,
		EntryType: "IT",
	})
	require.NoError(t, err)

	transactions, err := readTransactionFileData()
	require.NoError(t, err)
	require.NotEmpty(t, transactions)

	test_prams := []struct {
		name         string
		accountId    int64
		transactions []db.Transaction
	}{
		{"positive case", account.ID, transactions},
		{"account balance not enough", account2.ID, transactions},
		{"buy transaction cost not match", account.ID, transactions[len(transactions)-1:]},
		{"investment volume is lesser than transaction sell volume", account.ID, transactions[len(transactions)-2 : len(transactions)-1]},
		{"sell transaction return not match", account3.ID, transactions},
	}

	for _, tc := range test_prams {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.name {
			case "positive case":
				importedTransactions, err := tranService.ImportTransaction(context.Background(), tc.accountId, tc.transactions)
				require.NoError(t, err)
				require.NotNil(t, importedTransactions)
				for i, importedTrans := range importedTransactions {

					transacion := tc.transactions[i]
					require.Equal(t, transacion.Ticker, importedTrans.Ticker)
					require.Equal(t, transacion.TradingDate.Time.UTC(), importedTrans.TradingDate.Time)
					require.Equal(t, transacion.Trade, importedTrans.Trade)
					require.Equal(t, transacion.Volume, importedTrans.Volume)
					require.Equal(t, transacion.OrderPrice, importedTrans.OrderPrice)
					require.Equal(t, transacion.MatchVolume, importedTrans.MatchVolume)
					require.Equal(t, transacion.MatchPrice, importedTrans.MatchPrice)
					require.Equal(t, transacion.MatchValue, importedTrans.MatchValue)
					require.Equal(t, transacion.Fee, importedTrans.Fee)
					require.Equal(t, transacion.Tax, importedTrans.Tax)
					require.Equal(t, transacion.Cost, importedTrans.Cost)
					require.Equal(t, transacion.Return, importedTrans.Return)
					require.Equal(t, transacion.Status, importedTrans.Status)
				}
				slices.Reverse(tc.transactions)
			case "account balance not enough":
				_, err := tranService.ImportTransaction(context.Background(), tc.accountId, tc.transactions)
				require.Error(t, err)
				require.EqualError(t, err, "account balance is less than transation cost")
				slices.Reverse(tc.transactions)
			case "buy transaction cost not match":
				originTxCost := tc.transactions[0].Cost
				tc.transactions[0].Cost = 0
				_, err := tranService.ImportTransaction(context.Background(), tc.accountId, tc.transactions)
				require.EqualError(t, err, fmt.Sprintf("%s - %s transaction's cost is not match with capital cost in investment", tc.transactions[0].Ticker, tc.transactions[0].Trade))
				tc.transactions[0].Cost = originTxCost
			case "investment volume is lesser than transaction sell volume":
				_, err := tranService.ImportTransaction(context.Background(), tc.accountId, tc.transactions)
				require.EqualError(t, err, fmt.Sprintf("%s - %s transaction's match volume is lesser than investment volume", tc.transactions[0].Ticker, tc.transactions[0].Trade))
			case "sell transaction return not match":
				var manipulatedTransaction = tc.transactions[len(transactions)-2]
				tc.transactions[len(transactions)-2].Return = -1
				_, err := tranService.ImportTransaction(context.Background(), tc.accountId, tc.transactions)
				require.EqualError(t, err, fmt.Sprintf("%s - %s transaction's return value is not match", manipulatedTransaction.Ticker, manipulatedTransaction.Trade))
			}

		})
	}
}

func readTransactionFileData() ([]db.Transaction, error) {
	excelFile, err := os.Open("./data/105CA35050.xlsx")
	if err != nil {
		return nil, err
	}
	f, err := excelize.OpenReader(excelFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
	}()
	rows, err := f.GetRows("Sheet 1")
	if err != nil {
		return nil, err
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
			volume, _ := strconv.Atoi(row[3])
			orderPrice, _ := strconv.Atoi(strings.Replace(row[4], ",", "", -1))
			matchVolume, _ := strconv.Atoi(row[5])
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
	return transactions, nil
}

func TestGetSumarizeTransactionInfo(t *testing.T) {
	user := createRandomUser(t)
	acc := createRandomAccount(t, user.Username)

	transactions, err := readTransactionFileData()
	require.NoError(t, err)
	require.NotEmpty(t, transactions)

	result, err := tranService.GetSummarizeInfo(context.Background(), db.GetTransactionSummarizeInfoParams{
		AccountIds: []int64{acc.ID},
		Ticker:     transactions[0].Ticker,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, result.TotalRows, int32(0))
	require.GreaterOrEqual(t, result.SumFee, int64(0))
	require.GreaterOrEqual(t, result.SumMatchValue, int64(0))
	require.GreaterOrEqual(t, result.SumReturn, int64(0))
	require.GreaterOrEqual(t, result.SumTax, int64(0))

	importedTransactions, err := tranService.ImportTransaction(context.Background(), acc.ID, transactions)
	require.NoError(t, err)
	require.NotEmpty(t, importedTransactions)

	result, err = tranService.GetSummarizeInfo(context.Background(), db.GetTransactionSummarizeInfoParams{
		AccountIds: []int64{acc.ID},
		Ticker:     "",
	})

	var sumFee int64
	var sumMatchValue int64
	var sumReturn int64
	var sumTax int64
	for _, tx := range importedTransactions {
		sumFee += tx.Fee
		sumMatchValue += tx.MatchValue
		sumTax += tx.Tax
		sumReturn += tx.Return
	}
	require.Equal(t, int32(len(importedTransactions)), result.TotalRows)
	require.Equal(t, sumFee, result.SumFee)
	require.Equal(t, sumMatchValue, result.SumMatchValue)
	require.Equal(t, sumTax, result.SumTax)
	require.Equal(t, sumReturn, result.SumReturn)
}
