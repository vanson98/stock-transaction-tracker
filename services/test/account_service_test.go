package service_test

import (
	"context"
	"fmt"
	"strings"
	db "stt/database/postgres/sqlc"
	"stt/services"
	"stt/services/dtos"
	"stt/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T, userName string) db.Account {
	arg := db.CreateAccountParams{
		ChannelName: strings.ToUpper(util.RandomString(3)),
		Owner:       userName,
		Balance:     util.RandomInt(500000000, 1000000000),
		Currency:    util.RandomCurrency(),
	}
	account, err := accService.CreateNew(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.ChannelName, account.ChannelName)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)
	createRandomAccount(t, user.Username)
}

func TestGetById(t *testing.T) {
	user := createRandomUser(t)
	acc := createRandomAccount(t, user.Username)
	getAcc, err := accService.GetById(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Equal(t, acc, getAcc)
}

func TestGetStockAssetOverview(t *testing.T) {
	//create random account
	user := createRandomUser(t)
	acc := createRandomAccount(t, user.Username)
	// create a  investment
	ivm := createDefaultInvestmnet(t, acc.ID)
	// create transactions
	buyVolume := util.RandomInt(500, 1000)
	createRandomBuyTransaction(t, acc.ID, ivm.ID, buyVolume)
	sellVolume := util.RandomInt(100, 500)
	createRandomSellTransaction(t, acc.ID, ivm.ID, sellVolume)

	// check asset overview
	ivmDb, err := investmentService.GetById(context.Background(), ivm.ID)
	require.NoError(t, err)
	updateAccountBalance := acc.Balance - ivmDb.BuyValue - ivmDb.Fee - ivmDb.Tax + ivmDb.SellValue

	stockAssetOverview, err := accService.GetStockAssetOverview(context.Background(), []int64{acc.ID})
	require.NoError(t, err)
	require.Equal(t, acc.ID, stockAssetOverview[0].ID)
	require.Equal(t, acc.ChannelName, stockAssetOverview[0].ChannelName)
	require.Equal(t, updateAccountBalance, stockAssetOverview[0].Cash)
	require.Equal(t, ivmDb.CurrentVolume*ivmDb.CapitalCost, stockAssetOverview[0].TotalCogs)
	require.Equal(t, ivmDb.CurrentVolume*ivmDb.MarketPrice, stockAssetOverview[0].MarketValue)
}

func TestListAllByOwner(t *testing.T) {
	user := createRandomUser(t)
	n := 5
	for i := 0; i < n; i++ {
		createRandomAccount(t, user.Username)
	}
	accounts, err := accService.ListAllByOwner(context.Background(), user.Username)
	require.NoError(t, err)
	require.Greater(t, len(accounts), 0)
	require.Equal(t, n, len(accounts))
}

func TestGetAllOverview(t *testing.T) {
	//create random accounts
	user := createRandomUser(t)
	account := createRandomAccount(t, user.Username)
	var totalDeposit int64 = 0
	var totalWithdraw int64 = 0
	for i := 0; i < 5; i++ {
		randomAmount := util.RandomInt(-10000000, 10000000)

		result, err := accService.TransferMoney(context.Background(), dtos.TransferMoneyTxParam{
			AccountID: account.ID,
			Amount:    randomAmount,
			EntryType: db.EntryTypeIT,
		})
		require.NoError(t, err)
		require.NotEmpty(t, result)
		if result.Entry.Amount > 0 {
			totalDeposit += result.Entry.Amount
		} else {
			totalWithdraw += result.Entry.Amount
		}
	}

	// get all accounts by owner
	accountOverviews, err := accService.GetAllOverview(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, account.Balance+totalDeposit+totalWithdraw, accountOverviews[0].Balance)
	require.Equal(t, account.ChannelName, accountOverviews[0].ChannelName)
	require.Equal(t, account.ID, accountOverviews[0].ID)
	require.Equal(t, account.Currency, accountOverviews[0].Currency)
	require.Equal(t, totalDeposit, accountOverviews[0].Deposit)
	require.Equal(t, totalWithdraw, accountOverviews[0].Withdraw)
}

func TestTranserMoneyTx(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user.Username)
	fmt.Println(">> before:", account.Balance)
	n := 5
	var amount int64 = 10

	errChan := make(chan error)
	resultChan := make(chan dtos.TransferMoneyTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i)
		go func() {
			ctx := context.WithValue(context.Background(), services.TxKey, txName)
			result, err := accService.TransferMoney(ctx, dtos.TransferMoneyTxParam{
				AccountID: account.ID,
				Amount:    amount,
				EntryType: db.EntryTypeIT,
			})
			errChan <- err
			resultChan <- result
		}()
	}
	var existed map[int]bool = make(map[int]bool, 0)
	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)

		result := <-resultChan
		require.NotEmpty(t, result)

		entry := result.Entry
		updatedAccount := result.UpdatedAccount

		// check entry
		require.NotEmpty(t, entry)
		require.Equal(t, entry.Amount, amount)
		require.Equal(t, entry.AccountID, updatedAccount.ID)
		require.NotZero(t, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)
		require.NotZero(t, entry.Type)
		require.NotZero(t, entry.AccountID)

		_, err = store.GetEntryById(context.Background(), entry.ID)
		require.NoError(t, err)

		// check account
		require.NotEmpty(t, updatedAccount)
		require.Equal(t, updatedAccount.ID, account.ID)
		require.Equal(t, updatedAccount.Owner, account.Owner)
		require.Equal(t, updatedAccount.ChannelName, account.ChannelName)
		require.Equal(t, updatedAccount.CreatedAt, account.CreatedAt)

		// check balance
		fmt.Println(">> tx:", updatedAccount.Balance)
		diff := updatedAccount.Balance - account.Balance
		require.True(t, diff%amount == 0)
		require.Equal(t, account.Balance+diff, updatedAccount.Balance)
		k := int(diff / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
}
