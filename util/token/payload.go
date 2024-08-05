package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	minSecretKeySize = 32
)

var (
	ErrInvalidSecretKey = fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	ErrInvalidToken     = fmt.Errorf("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
)

// Payload is a struct to represent the payload of a token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Sub       string    `json:"sub"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(sub string, duration time.Duration) (*Payload, error) {

	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Sub:       sub,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
