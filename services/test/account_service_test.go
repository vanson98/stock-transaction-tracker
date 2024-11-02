package service_test

import (
	"context"
	"fmt"
	db "stt/database/postgres/sqlc"
	"stt/services"
	"stt/services/dtos"
	"stt/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) db.Account {
	user := createRandomUser(t)
	arg := db.CreateAccountParams{
		ChannelName: util.RandomString(3),
		Owner:       user.Username,
		Balance:     util.RandomInt(1, 1000),
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
	createRandomAccount(t)
}

func TestGetById(t *testing.T) {
	acc := createRandomAccount(t)
	getAcc, err := accService.GetById(context.Background(), acc.ID)
	require.NoError(t, err)
	require.Equal(t, acc, getAcc)
}

func TestGetAllPaging(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomAccount(t)
	}
	accounts, err := accService.GetAllPaging(context.Background(), db.GetAccountsPagingParams{
		Limit:  5,
		Offset: 5,
	})
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, a := range accounts {
		require.NotEmpty(t, a)
	}
}

func TestUpdateAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t)
	param := db.AddAccountBalanceParams{
		ID:     account1.ID,
		Amount: util.RandomInt(1, 100),
	}

	account2, err := accService.UpdateBalance(context.Background(), param)
	require.NoError(t, err)
	require.NotNil(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance+param.Amount, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	err := accService.Delete(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := store.GetAccountById(context.Background(), acc1.ID)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestTranserMoneyTx(t *testing.T) {
	account := createRandomAccount(t)
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
	accounts, err := accService.ListAllAccount(context.Background())
	require.NoError(t, err)
	require.Greater(t, len(accounts), 0)
}
