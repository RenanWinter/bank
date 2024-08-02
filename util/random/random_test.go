package random

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInt(t *testing.T) {
	min := int64(10)
	max := int64(20)
	result := Int(min, max)
	require.GreaterOrEqual(t, result, min)
	require.LessOrEqual(t, result, max)
}

func TestString(t *testing.T) {
	result := String(10)
	require.Equal(t, 10, len(result))
}

func TestEmail(t *testing.T) {
	result := Email()
	require.Contains(t, result, "@")
	require.Contains(t, result, ".com")
	require.Equal(t, 15, len(result))
}

func TestDateTime(t *testing.T) {
	result := DateTime(nil, nil)
	dummy := time.Now()
	require.NotZero(t, result)
	require.IsType(t, dummy, result)
}
func TestDateTimeWithMinAndMax(t *testing.T) {
	min := time.Now()
	max := time.Now().Add(time.Hour * 48)
	result := DateTime(&min, &max)
	dummy := time.Now()
	require.NotZero(t, result)
	require.IsType(t, dummy, result)
	require.True(t, result.Before(max))
	require.True(t, result.After(min))
}

func TestDecimal(t *testing.T) {
	min := float64(10)
	max := float64(10.99)
	places := 2
	result := Decimal(min, max, places)
	require.GreaterOrEqual(t, result, min)
	require.LessOrEqual(t, result, max)
}
