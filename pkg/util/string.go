package util

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomString generates a random string of length n
func RandomString(n int) string {
	id := gonanoid.MustGenerate(letterBytes, n)
	return id
}
