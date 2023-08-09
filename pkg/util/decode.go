package util

import (
	"encoding/base64"
	"fmt"
)

// ParseEnvSecret load firebase admin token secret from env
func ParseEnvSecret(encodedToken string) ([]byte, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return nil, fmt.Errorf("error decoding Firebase Admin SDK token: %v", err)
	}

	return decodedToken, nil
}
