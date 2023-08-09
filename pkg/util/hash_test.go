package util

import (
	"testing"
)

func TestGenerateHashedKey(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "abcd1234",
			expected: "$2a$10$0qtOvmMkMEYUlfoYy/8/..sdAaKVbiYhJRvRQ0wVMf54Gwkys4oWW", // Error case; don't check the actual hash value, as it changes every time
		},
	}

	for _, test := range tests {
		hashedKey, err := GenerateHashedKey(test.input)

		if test.expected == "" {
			// Error case
			if err == nil {
				t.Errorf("GenerateHashedKey(%s) expected an error, but got nil", test.input)
			}
		} else {
			// Success case
			if err != nil {
				t.Errorf("GenerateHashedKey(%s) returned an error: %v", test.input, err)
			}
			if hashedKey == "" {
				t.Errorf("GenerateHashedKey(%s) returned an empty hashed key", test.input)
			}
		}
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		hashed   string
		expected bool
	}{
		{
			password: "abcd1234",
			hashed:   "$2a$10$0qtOvmMkMEYUlfoYy/8/..sdAaKVbiYhJRvRQ0wVMf54Gwkys4oWW",
			expected: true,
		},
		{
			password: "wrongpassword",
			hashed:   "$2a$10$0qtOvmMkMEYUlfoYy/8/..sdAaKVbiYhJRvRQ0wVMf54Gwkys4oWW",
			expected: false,
		},
	}

	for _, test := range tests {
		actual := IsValidPassword(test.password, test.hashed)
		if actual != test.expected {
			t.Errorf("IsValidPassword(%s, %s) expected %t, but got %t", test.password, test.hashed, test.expected, actual)
		}
	}
}
