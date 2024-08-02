package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RenanWinter/bank/util/random"
)

var account Account

func _createFakeAccount(t *testing.T, user User) Account {
	arg := CreateAccountParams{
		Name:          random.String(10),
		OwnerID:       user.ID,
		AccountTypeID: 1,
		Balance:       0,
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)
	require.Equal(t, arg.Name, account.Name)
	require.Equal(t, arg.OwnerID, account.OwnerID)
	require.Equal(t, arg.AccountTypeID, account.AccountTypeID)
	require.Equal(t, arg.Balance, account.Balance)
	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.UpdatedAt)
	require.Zero(t, account.DeletedAt)

	return account
}

func _getAccount(t *testing.T) Account {
	if account.ID != 0 {
		return account
	}

	user := _getUser(t)
	account = _createFakeAccount(t, user)
	return account
}

func TestCreateAccount(t *testing.T) {
	_getAccount(t)
}
