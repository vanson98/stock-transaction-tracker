package service_test

import (
	"context"
	"fmt"
	db "stt/database/postgres/sqlc"
	"stt/services"
	"stt/services/dtos"
	"stt/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T, user db.User) db.Account {
	arg := db.CreateAccountParams{
		ChannelName: util.RandomString(3),
		Owner:       user.Username,
		Balance:     50000000,
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
	createRandomAccount(t, user)
}

func TestGetById(t *testing.T) {
	user := createRandomUser(t)
	acc := createRandomAccount(t, user)
	getAcc, err := accService.GetById(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Equal(t, acc, getAcc)
}

func TestTranserMoneyTx(t *testing.T) {
	user := createRandomUser(t)
	account := createRandomAccount(t, user)
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

func TestListAllAccount(t *testing.T) {
	accounts, err := accService.ListAllByOwner(context.Background(), "vanson")
	require.NoError(t, err)
	require.Greater(t, len(accounts), 0)
}

func TestGetAccountInfoByIds(t *testing.T) {
	// create random account
	user := createRandomUser(t)
	acc := createRandomAccount(t, user)
	// create a depositTransfer money depositTransfer
	depositTransfer, err := accService.TransferMoney(context.Background(), dtos.TransferMoneyTxParam{
		AccountID: acc.ID,
		Amount:    util.RandomInt(50, 100),
		EntryType: db.EntryTypeTM,
	})

	require.NoError(t, err)
	require.Equal(t, acc.ID, depositTransfer.Entry.AccountID)
	require.Equal(t, depositTransfer.Entry.Type, db.EntryTypeTM)

	// create a withdrawal money entry
	withdrawalTransfer, err := accService.TransferMoney(context.Background(), dtos.TransferMoneyTxParam{
		AccountID: acc.ID,
		Amount:    util.RandomInt(-50, 0),
		EntryType: db.EntryTypeTM,
	})

	require.NoError(t, err)
	require.Equal(t, withdrawalTransfer.Entry.AccountID, acc.ID)
	require.Equal(t, withdrawalTransfer.Entry.Type, db.EntryTypeTM)

	// get account info
	accInfo, err := accService.GetAccountInfoByIds(context.Background(), []int64{acc.ID})
	require.NoError(t, err)
	require.Equal(t, accInfo[0].ID, acc.ID)
	require.Equal(t, accInfo[0].Cash, withdrawalTransfer.UpdatedAccount.Balance)
}

func TestGetAllByOwner(t *testing.T) {
	// create random accounts
	// accounts, username := createRandomListAccountForUser(t)

	// // get all accounts by owner
	// allAccounts, err := accService.GetAllByOwner(context.Background(), "vanson")
	// require.NoError(t, err)
	// require.Equal(t, len(accounts), len(allAccounts))
	// for i, acc := range allAccounts {
	// 	require.Equal(t, acc.ID, accounts[i].ID)
	// 	require.Equal(t, acc.Owner, accounts[i].Balance)

	// }
}
