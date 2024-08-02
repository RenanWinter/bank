package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveUser(t *testing.T) {
	user := _createFakeUser(t)
	_createFakeCredential(t, user)
	_createFakeCredential(t, user)
	store := NewStore(testDB)
	err := store.RemoveUser(context.Background(), user.ID)
	require.NoError(t, err)
	stored, err := testQueries.GetUserById(context.Background(), user.ID)
	require.Error(t, err)
	require.Empty(t, stored)

	credential, err := testQueries.GetUserActiveCredential(context.Background(), user.ID)
	require.Error(t, err)
	require.Empty(t, credential)

}
