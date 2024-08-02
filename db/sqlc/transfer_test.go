package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfer(t *testing.T) {
	store := NewStore(testDB)
	user2 := _createFakeUser(t)
	fromAccount := _getAccount(t)
	toAccount := _createFakeAccount(t, user2)

	// run n concurrent transfser transactions
	n := 10
	amount := float64(10)

	// create a channel to receive the results and errors
	results := make(chan TransferResult)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.Transfer(ctx, TransferParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result

		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check the transfer
		tranfer := result.Transfer
		require.NotEmpty(t, tranfer)
		require.Equal(t, fromAccount.ID, tranfer.FromAccountID)
		require.Equal(t, toAccount.ID, tranfer.ToAccountID)
		require.Equal(t, amount, tranfer.Amount)
		require.NotZero(t, tranfer.ID)
		require.NotZero(t, tranfer.CreatedAt)
		require.NotZero(t, tranfer.UpdatedAt)
		require.Zero(t, tranfer.DeletedAt)

		_, err = store.GetTransfer(context.Background(), tranfer.ID)
		require.NoError(t, err)

		//check withdraw movement
		fromMovement := result.FromMovement
		require.NotEmpty(t, fromMovement)
		require.Equal(t, fromAccount.ID, fromMovement.AccountID)
		require.Equal(t, -amount, fromMovement.Amount)
		require.Equal(t, tranfer.ID, fromMovement.TransferID.Int64)
		require.Equal(t, "Transfer to account "+toAccount.Name, fromMovement.Description)
		require.False(t, fromMovement.Validated)
		require.NotZero(t, fromMovement.ID)
		require.NotZero(t, fromMovement.CreatedAt)
		require.NotZero(t, fromMovement.UpdatedAt)
		require.Zero(t, fromMovement.DeletedAt)

		_, err = store.GetMovement(context.Background(), fromMovement.ID)
		require.NoError(t, err)

		//check deposit movement
		toMovement := result.ToMovement
		require.NotEmpty(t, toMovement)
		require.Equal(t, toAccount.ID, toMovement.AccountID)
		require.Equal(t, amount, toMovement.Amount)
		require.Equal(t, tranfer.ID, toMovement.TransferID.Int64)
		require.Equal(t, "Transfer from account "+fromAccount.Name, toMovement.Description)
		require.False(t, toMovement.Validated)
		require.NotZero(t, toMovement.ID)
		require.NotZero(t, toMovement.CreatedAt)
		require.NotZero(t, toMovement.UpdatedAt)
		require.Zero(t, toMovement.DeletedAt)

		_, err = store.GetMovement(context.Background(), toMovement.ID)
		require.NoError(t, err)

		// check the accounts balance
		_fromAccount := result.FromAccount
		require.NotEmpty(t, _fromAccount)
		require.Equal(t, _fromAccount.ID, fromAccount.ID)

		_toAccount := result.ToAccount
		require.NotEmpty(t, _toAccount)
		require.Equal(t, _toAccount.ID, toAccount.ID)

		diff1 := fromAccount.Balance - _fromAccount.Balance
		diff2 := _toAccount.Balance - toAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check final updatedBalance
	updatedFromAccount, err := store.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	require.Equal(t, fromAccount.Balance-float64(n)*amount, updatedFromAccount.Balance)

	updatedToAccount, err := store.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)
	require.Equal(t, toAccount.Balance+float64(n)*amount, updatedToAccount.Balance)
}

func TestTransferDeadlock(t *testing.T) {
	store := NewStore(testDB)
	user1 := _createFakeUser(t)
	user2 := _createFakeUser(t)

	account1 := _createFakeAccount(t, user1)
	account2 := _createFakeAccount(t, user2)

	n := 10
	amount := float64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		go func() {
			_, err := store.Transfer(context.Background(), TransferParams{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			fmt.Println("Transfer error:", i+1, err)
		}
		require.NoError(t, err)
	}

	// check final updatedBalance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
