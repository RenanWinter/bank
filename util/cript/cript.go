package cript

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, cost int) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		if err == bcrypt.ErrPasswordTooLong {
			return "", fmt.Errorf("Your password is too long. Try again with a password with less than 72 characters.")
		}
		return "", fmt.Errorf("Failed to protect your account. Try again in some minutes.")
	}
	return string(hashed), nil
}

func CheckPassword(password, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
