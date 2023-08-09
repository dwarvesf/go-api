package util

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// GenerateHashedKey hash key
func GenerateHashedKey(key string) (string, error) {
	val := strings.TrimSpace(key)
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(val), bcrypt.DefaultCost)
	return string(hashedKey), err
}

// IsValidPassword validate the password
func IsValidPassword(val string, expected string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(expected), []byte(val))
	return err == nil
}
