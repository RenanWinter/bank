package random

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

var randomizer *rand.Rand

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Int(min, max int64) int64 {
	return min + randomizer.Int63n(max-min+1)
}

func String(n int) string {
	var sb strings.Builder

	k := len(letters)

	for i := 0; i < n; i++ {
		c := letters[randomizer.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func Email() string {
	return String(6) + "@" + String(4) + ".com"
}

func DateTime(min, max *time.Time) time.Time {
	if min == nil {
		min = new(time.Time)
	}

	if max == nil {
		max = new(time.Time)
		*max = time.Now()
	}

	delta := max.Sub(*min)
	offset := time.Duration(randomizer.Int63n(int64(delta)))

	return min.Add(offset)
}

func Decimal(min, max float64, places int) float64 {
	value := min + randomizer.Float64()*(max-min)
	precision := 1.0

	for i := 0; i < places; i++ {
		precision *= 10
	}

	return float64(int(value*precision)) / precision
}

func Bool() bool {
	return randomizer.Intn(2) == 1
}

func Duration(min, max time.Duration) time.Duration {
	return time.Duration(Int(int64(min), int64(max)))
}

func UUID() string {
	return uuid.New().String()
}
