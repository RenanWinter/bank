package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWTMaker is a struct JSON Web Token Maker that implements the Maker interface
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, ErrInvalidSecretKey
	}

	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(sub string, duration time.Duration) (string, error) {
	payload, err := NewPayload(sub, duration)
	if err != nil {
		return "", err
	}

	jwtTokken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtTokken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
