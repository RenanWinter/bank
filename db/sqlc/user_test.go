package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RenanWinter/bank/util/random"
)

var _user User

func _getUser(t *testing.T) User {
	if _user.ID != 0 {
		return _user
	}

	_user = _createFakeUser(t)
	return _user
}

func _createFakeUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: random.String(10),
		Email:    random.Email(),
		Name:     fmt.Sprintf("%v %v", random.String(6), random.String(6)),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.Uuid)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Name, user.Name)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)
	require.Zero(t, user.DeletedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	_getUser(t)
}

func TestUsersUniqueConstraints(t *testing.T) {
	user := _getUser(t)

	arg := CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
	}

	_, err := testQueries.CreateUser(context.Background(), arg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate key value violates unique constraint")
}

func TestGetUserByEmail(t *testing.T) {
	user := _createFakeUser(t)

	user2, err := testQueries.GetUserByEmail(context.Background(), user.Email)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Uuid, user2.Uuid)
	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.CreatedAt, user2.CreatedAt)
	require.Equal(t, user.UpdatedAt, user2.UpdatedAt)
	require.Equal(t, user.DeletedAt, user2.DeletedAt)
}

func TestGetUserByEmailNotFound(t *testing.T) {
	_, err := testQueries.GetUserByEmail(context.Background(), random.Email())
	require.Error(t, err)
}

func TestGetUserById(t *testing.T) {
	user := _createFakeUser(t)

	user2, err := testQueries.GetUserById(context.Background(), user.ID)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Uuid, user2.Uuid)
	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.CreatedAt, user2.CreatedAt)
	require.Equal(t, user.UpdatedAt, user2.UpdatedAt)
	require.Equal(t, user.DeletedAt, user2.DeletedAt)
}

func TestGetUserByIdNotFound(t *testing.T) {
	_, err := testQueries.GetUserById(context.Background(), 99999999)
	require.Error(t, err)
}

func TestGetUserByUsername(t *testing.T) {
	user := _createFakeUser(t)

	user2, err := testQueries.GetUserByUsername(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Uuid, user2.Uuid)
	require.Equal(t, user.Username, user2.Username)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Name, user2.Name)
	require.Equal(t, user.CreatedAt, user2.CreatedAt)
	require.Equal(t, user.UpdatedAt, user2.UpdatedAt)
	require.Equal(t, user.DeletedAt, user2.DeletedAt)
}

func TestGetUserByUsernameNotFound(t *testing.T) {
	_, err := testQueries.GetUserByUsername(context.Background(), random.String(10))
	require.Error(t, err)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		_createFakeUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
