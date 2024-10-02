package db

import (
	"context"
	"stt/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranserTx(t *testing.T) {
	store := NewStore(pgConnPool)

	acc := createRandomAccount(t)
	n := 2
	amount := util.RandomInt(-50, 50)

	errChan := make(chan error)
	resultChan := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TranserTx(context.Background(), TransferTxParam{
				AccountID:     acc.ID,
				Amount:        amount,
				EntryFromType: EntryFromTypeIT,
			})
			errChan <- err
			resultChan <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errChan
		require.NoError(t, err)

		result := <-resultChan
		require.NotEmpty(t, result)

		entry := result.AccountEntry
		account := result.UpdatedAccount

		// check entry
		require.NotEmpty(t, entry)
		require.Equal(t, entry.Amount, amount)
		require.Equal(t, entry.AccountID, account.ID)
		require.NotZero(t, entry.AccountID)
		require.NotZero(t, entry.CreatedAt)
		require.NotZero(t, entry.FromType)
		require.NotZero(t, entry.AccountID)

		// check entry exist
		_, err = testQueries.GetAccountEntryById(context.Background(), entry.ID)
		require.NoError(t, err)

		// check balance
		//require.Equal(t, acc.Balance+amount, account.Balance)

	}

}
