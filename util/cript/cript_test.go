package cript

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/RenanWinter/bank/util/random"
)

func TestHashPassword(t *testing.T) {
	password := random.String(20)

	hashed1, err := HashPassword(password, bcrypt.DefaultCost)
	require.NoError(t, err)
	require.NotEmpty(t, hashed1)

	err = CheckPassword(password, hashed1)
	require.NoError(t, err)

	wrongPassword := random.String(20)

	err = CheckPassword(wrongPassword, hashed1)
	require.Error(t, err)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashed2, err := HashPassword(password, bcrypt.DefaultCost)
	require.NoError(t, err)
	require.NotEmpty(t, hashed2)
	require.NotEqual(t, hashed1, hashed2)

	longPassword := random.String(100)
	hashed3, err := HashPassword(longPassword, bcrypt.DefaultCost)
	require.Error(t, err)
	require.Empty(t, hashed3)
	require.EqualError(t, err, "Your password is too long. Try again with a password with less than 72 characters.")

	minCost := bcrypt.MinCost - 1
	hashed4, err := HashPassword(password, minCost)
	require.NoError(t, err)
	require.NotEmpty(t, hashed4)

	maxCost := bcrypt.MaxCost + 1
	hashed5, err := HashPassword(password, maxCost)
	require.Error(t, err)
	require.Empty(t, hashed5)
	require.EqualError(t, err, "Failed to protect your account. Try again in some minutes.")

}

func TestCheckPassword(t *testing.T) {
	password := random.String(20)

	hashed, err := HashPassword(password, bcrypt.DefaultCost)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	err = CheckPassword(password, hashed)
	require.NoError(t, err)

	wrongPassword := random.String(20)

	err = CheckPassword(wrongPassword, hashed)
	require.Error(t, err)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
