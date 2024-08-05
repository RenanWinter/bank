package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/RenanWinter/bank/util/random"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(random.String(32))
	require.NoError(t, err)

	uuid := random.UUID()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(uuid, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, uuid, payload.Sub)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(random.String(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(random.UUID(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
