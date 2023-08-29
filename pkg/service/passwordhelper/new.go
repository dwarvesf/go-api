package passwordhelper

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"io"
)

const (
	// saltSize is the size of the salt used to hash passwords
	saltSize = 32
)

// Helper password helper
type Helper interface {
	GenerateSalt() string
	Hash(password, salt string) (string, error)
	Compare(password, hashedPassword, salt string) bool
}

type implSha512 struct{}

// NewHelper init helper
func NewHelper() Helper {
	return &implSha512{}
}

func (h implSha512) GenerateSalt() string {
	var salt = make([]byte, saltSize)

	_, err := io.ReadFull(rand.Reader, salt)

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(salt)

}

func (h implSha512) Hash(password, salt string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)
	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		return "", err
	}

	// Create sha-512 hasher
	var sha512Hasher = sha512.New()
	// Append salt to password
	passwordBytes = append(passwordBytes, saltBytes...)

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	// Convert the hashed password to a hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)

	// Return the hashed password
	return hashedPasswordHex, nil
}

func (h implSha512) Compare(password, hashedPassword, salt string) bool {
	currPasswordHash, err := h.Hash(password, salt)
	if err != nil {
		return false
	}

	return hashedPassword == currPasswordHash
}
