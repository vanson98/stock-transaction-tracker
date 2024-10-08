package db

import (
	"context"
	"stt/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		ChannelName: util.RandomString(3),
		Owner:       util.RandomOwner(),
		Balance:     util.RandomInt(1, 1000),
		Currency:    util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
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

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	accounts, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
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
	param := AddAccountBalanceParams{
		ID:     account1.ID,
		Amount: util.RandomInt(1, 100),
	}

	account2, err := testQueries.AddAccountBalance(context.Background(), param)
	require.NoError(t, err)
	require.NotNil(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance+param.Amount, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}
