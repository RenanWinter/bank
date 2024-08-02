package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RenanWinter/bank/util/cript"
	"github.com/RenanWinter/bank/util/random"
)

var _credential Credential

func _createFakeCredential(t *testing.T, user User) Credential {
	password := cript.HashPassword(random.String(20))

	arg := CreateCredentialParams{
		UserID:   user.ID,
		Password: password,
	}

	credential, err := testQueries.CreateCredential(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, credential)
	require.Equal(t, arg.UserID, credential.UserID)
	require.Equal(t, arg.Password, credential.Password)
	require.NotZero(t, credential.ID)
	require.NotZero(t, credential.CreatedAt)
	require.NotZero(t, credential.UpdatedAt)
	require.Zero(t, credential.DeletedAt)
	require.Equal(t, 60, len(credential.Password))
	return credential

}
func _getCredential(t *testing.T) Credential {

	if _credential.ID != 0 {
		return _credential
	}

	user := _getUser(t)
	_credential = _createFakeCredential(t, user)
	return _credential
}

func TestCreateCredential(t *testing.T) {
	_getCredential(t)
}

func TestGetUserActiveCredential(t *testing.T) {
	newCredential := _getCredential(t)
	credential2, err := testQueries.GetUserActiveCredential(context.Background(), newCredential.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, credential2)
	require.Equal(t, newCredential.ID, credential2.ID)
	require.Equal(t, newCredential.UserID, credential2.UserID)
	require.Equal(t, newCredential.Password, credential2.Password)
	require.Equal(t, newCredential.CreatedAt, credential2.CreatedAt)
	require.Equal(t, newCredential.UpdatedAt, credential2.UpdatedAt)
	require.Equal(t, newCredential.DeletedAt, credential2.DeletedAt)
}

func TestGetUserActiveCredentialNoRows(t *testing.T) {
	_, err := testQueries.GetUserActiveCredential(context.Background(), 9999999)
	require.Error(t, err)
}

func TestGetUserCredentials(t *testing.T) {
	user := _createFakeUser(t)
	_createFakeCredential(t, user)
	_createFakeCredential(t, user)
	credentials, err := testQueries.GetUserCredentials(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, credentials)
	require.Equal(t, len(credentials), 2)
}

func TestGetUserCredentialsNoRows(t *testing.T) {
	credentials, err := testQueries.GetUserCredentials(context.Background(), 9999999)
	require.NoError(t, err)
	require.Empty(t, credentials)
}

func TestRemoveUserCredential(t *testing.T) {
	credential := _getCredential(t)
	err := testQueries.RemoveUserCredential(context.Background(), credential.ID)
	require.NoError(t, err)
}
