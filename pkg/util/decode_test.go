package util

import (
	"bytes"
	"testing"
)

func TestParseEnvSecret(t *testing.T) {
	testCases := map[string]struct {
		encodedToken string
		wantValue    []byte
		wantErr      bool
	}{
		"Valid encoded token": {
			encodedToken: "cGFzc3dvcmQ=",
			wantValue:    []byte("password"),
			wantErr:      false,
		},
		"Invalid encoded token": {
			encodedToken: "not_base64_encoded",
			wantValue:    nil,
			wantErr:      true,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
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
