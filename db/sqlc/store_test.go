package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testPool)

	sender := createRandomAccount(t)
	receiver := createRandomAccount(t)

	// channel
	errs := make(chan error)
	res := make(chan TransferTxResult)

	// test for concurrency
	n := 5
	amount := 10
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx-%d", i+1)
		fmt.Println(txName)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: sender.ID,
				ToAccountID:   receiver.ID,
				Amount:        float64(amount),
			})

			errs <- err
			res <- result
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		r := <-res
		require.NotEmpty(t, r)

		// check transfer obj
		transfer := r.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, sender.ID)
		require.Equal(t, transfer.ToAccountID, receiver.ID)
		require.Equal(t, transfer.Amount, float64(amount))
		require.NotEmpty(t, transfer.ID)
		require.NotEmpty(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from entry
		fromEntry := r.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotEmpty(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry.CreatedAt)
		require.Equal(t, fromEntry.AccountID, sender.ID)
		require.Equal(t, fromEntry.Amount, float64(-amount))

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check to entry
		toEntry := r.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotEmpty(t, toEntry.ID)
		require.NotEmpty(t, toEntry.CreatedAt)
		require.Equal(t, toEntry.AccountID, receiver.ID)
		require.Equal(t, toEntry.Amount, float64(amount))

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check from account
		fromAccount := r.FromAccount
		require.NotEmpty(t, r.FromAccount)
		require.Equal(t, fromAccount.ID, sender.ID)
		require.Equal(t, fromAccount.Currency, sender.Currency)
		require.Equal(t, fromAccount.Owner, sender.Owner)

		// check to account
		toAccount := r.ToAccount
		require.NotEmpty(t, r.ToAccount)
		require.Equal(t, toAccount.ID, receiver.ID)
		require.Equal(t, toAccount.Currency, receiver.Currency)
		require.Equal(t, toAccount.Owner, receiver.Owner)

		fmt.Printf(">> tx: %f %f\n", fromAccount.Balance, toAccount.Balance)
		// check balance, can not send more money than in account
		diff1 := sender.Balance - fromAccount.Balance
		// check balance, must have revieved more money
		diff2 := toAccount.Balance - receiver.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2 > 0)
	}
	// check balance after concurrent transfers
	updatedSender, err := testQueries.GetAccount(context.Background(), sender.ID)
	require.NoError(t, err)

	updatedReciever, err := testQueries.GetAccount(context.Background(), receiver.ID)
	require.NoError(t, err)

	require.Equal(t, updatedSender.Balance, sender.Balance-float64(n)*float64(amount))
	require.Equal(t, updatedReciever.Balance, receiver.Balance+float64(n)*float64(amount))

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testPool)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// channel
	errs := make(chan error)

	// test for concurrency
	n := 10
	amount := 10
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx-%d", i+1)
		fmt.Println(txName)

		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        float64(amount),
			})

			errs <- err
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}
	// check balance after concurrent transfers
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, updatedAccount1.Balance, account1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)

}
