package passwordhelper

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/argon2"
)

type argon2Impl struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func newArgon2Default() *argon2Impl {
	return &argon2Impl{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

func (h argon2Impl) GenerateSalt() string {
	var salt = make([]byte, h.saltLength)

	_, err := io.ReadFull(rand.Reader, salt)

	if err != nil {
		panic(err)
	}

	return base64.RawStdEncoding.EncodeToString(salt)
}

func (h argon2Impl) Hash(password, salt string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(passwordBytes, saltBytes, h.iterations, h.memory, h.parallelism, h.keyLength)
	hashedStr := base64.RawStdEncoding.EncodeToString(hash)
	return hashedStr, nil
}

func (h argon2Impl) Compare(password, hashedPassword, salt string) bool {
	saltBytes, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return false
	}

	otherHash := argon2.IDKey([]byte(password), saltBytes, h.iterations, h.memory, h.parallelism, h.keyLength)

	hashedPasswordBytes, err := base64.RawStdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false
	}

	if subtle.ConstantTimeCompare(hashedPasswordBytes, otherHash) == 1 {
		return true
	}
	return false
}
