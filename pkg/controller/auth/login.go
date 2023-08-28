package auth

import (
	"time"

	"github.com/pkg/errors"

	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
	"github.com/golang-jwt/jwt/v5"
)

func (c *impl) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	// TODO: get user from db
	var user *orm.User

	now := time.Now()

	// Generate JWT token
	token, err := c.jwtHelper.GenerateJWTToken(map[string]interface{}{
		"sub": user.ID,
		"iss": c.cfg.App,
		"exp": jwt.NewNumericDate(now.AddDate(1, 0, 0)),
		"nbf": jwt.NewNumericDate(now),
		"iat": jwt.NewNumericDate(now),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &model.LoginResponse{
		ID:          user.ID,
		Email:       user.Email,
		AccessToken: token,
	}, nil
}
