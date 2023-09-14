package util

import (
	"strings"
	"testing"
)

func TestRandomString(t *testing.T) {
	testCases := map[string]struct {
		length int
	}{
		"TestLength1":   {1},
		"TestLength10":  {10},
		"TestLength100": {100},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			randomStr := RandomString(tc.length)
			if len(randomStr) != tc.length {
				t.Errorf("Expected string of length %d, got %d", tc.length, len(randomStr))
			}

			// Check if all characters in the random string are from letterBytes
			for _, char := range randomStr {
				if !strings.Contains(letterBytes, string(char)) {
					t.Errorf("Random string contains invalid character: %s", string(char))
				}
			}
		})
	}
}
