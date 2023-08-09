package util

import (
	"bytes"
	"testing"
)

func TestParseEnvSecret(t *testing.T) {
	testCases := []struct {
		name         string
		encodedToken string
		wantValue    []byte
		wantErr      bool
	}{
		{
			name:         "Valid encoded token",
			encodedToken: "cGFzc3dvcmQ=",
			wantValue:    []byte("password"),
			wantErr:      false,
		},
		{
			name:         "Invalid encoded token",
			encodedToken: "not_base64_encoded",
			wantValue:    nil,
			wantErr:      true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseEnvSecret(tt.encodedToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFirebaseSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !bytes.Equal(result, tt.wantValue) {
				t.Errorf("expected %v, but got %v", tt.wantValue, result)
			}
		})
	}
}
