package jwthelper

import (
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/golang-jwt/jwt/v5"
)

// Helper jwt helper
//
//go:generate mockery --name=Helper --with-expecter --output ./mocks
type Helper interface {
	GenerateJWTToken(claims jwt.MapClaims) (string, error)
	ValidateToken(token string) (map[string]interface{}, error)
}

type impl struct {
	Secret string
}

// NewHelper init helper
func NewHelper(secret string) Helper {
	return &impl{
		Secret: secret,
	}
}

func (h impl) GenerateJWTToken(claims jwt.MapClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(h.Secret))
}

func (h impl) ValidateToken(token string) (map[string]interface{}, error) {
	claims := &jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, model.ErrInvalidToken
	}
	return *claims, err
}
