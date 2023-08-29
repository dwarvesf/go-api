package jwthelper

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWTToken(t *testing.T) {
	// Test cases
	testCases := []struct {
		name      string
		secret    string
		claims    jwt.MapClaims
		expectErr bool
	}{
		{
			name:   "ValidToken",
			secret: "secret_key",
			claims: jwt.MapClaims{
				"sub": "user123",
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(time.Hour).Unix(),
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new instance of Helper with the test secret
			helper := NewHelper(tc.secret)

			// Generate the JWT token
			token, err := helper.GenerateJWTToken(tc.claims)

			// Check if an error is expected
			if tc.expectErr && err == nil {
				t.Errorf("Expected error, but got nil")
			} else if !tc.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// If there was no error, verify the token
			if err == nil {
				parsedToken, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					return []byte(tc.secret), nil
				})
				if parseErr != nil {
					t.Errorf("Failed to parse generated token: %v", parseErr)
				}

				// Check the validity of the claims
				if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
					for key, value := range tc.claims {
						switch tVal := value.(type) {
						case int64:
							newVal := claims[key].(float64)
							if tVal != int64(newVal) {
								t.Errorf("++++Claim mismatch for key %s: expected %v, got %v", key, value, claims[key])
							}

						default:
							if claims[key] != value {
								t.Errorf("Claim mismatch for key %s: expected %v, got %v", key, value, claims[key])
							}
						}
					}
				} else {
					t.Errorf("Failed to get claims from token")
				}
			}
		})
	}
}

func Test_impl_ValidateToken(t *testing.T) {
	h1 := &impl{
		Secret: "secret",
	}
	validTime := jwt.NewNumericDate(time.Now().AddDate(0, 0, 1))
	validToken, err := h1.GenerateJWTToken(jwt.MapClaims{
		"exp": validTime,
	})
	require.NoError(t, err)

	h2 := &impl{
		Secret: "secret invalid",
	}
	invalidToken, err := h2.GenerateJWTToken(jwt.MapClaims{
		"exp": validTime,
	})
	require.NoError(t, err)

	type fields struct {
		Secret string
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Secret: "secret",
			},
			args: args{
				token: validToken,
			},
			want: jwt.MapClaims{
				"exp": float64(validTime.Unix()),
			},
		},
		{
			name: "invalid token",
			fields: fields{
				Secret: "secret",
			},
			args: args{
				token: invalidToken,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := impl{
				Secret: tt.fields.Secret,
			}
			got, err := h.ValidateToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("impl.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("impl.ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
