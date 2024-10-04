package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranserTx(t *testing.T) {
	store := NewStore(pgConnPool)

	account := createRandomAccount(t)
	fmt.Println(">> before:", account.Balance)
	n := 100
	var amount int64 = 10

	errChan := make(chan error)
	resultChan := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TranserTx(ctx, TransferTxParam{
				AccountID: account.ID,
				Amount:    amount,
				EntryType: EntryTypeIT,
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

		entry := result.AccountEntry
		updatedAccount := result.UpdatedAccount

		// check entry
		require.NotEmpty(t, entry)
		require.Equal(t, entry.Amount, amount)
		require.Equal(t, entry.AccountID, updatedAccount.ID)
		require.NotZero(t, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)
		require.NotZero(t, entry.Type)
		require.NotZero(t, entry.AccountID)

		_, err = testQueries.GetEntryById(context.Background(), entry.ID)
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
