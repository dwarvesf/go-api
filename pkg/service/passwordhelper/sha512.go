package passwordhelper

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"io"
)

type implSha512 struct {
	saltLength uint32
}

func newSha512Default() *implSha512 {
	return &implSha512{
		saltLength: 32,
	}
}

func (h implSha512) GenerateSalt() string {
	var salt = make([]byte, h.saltLength)

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
