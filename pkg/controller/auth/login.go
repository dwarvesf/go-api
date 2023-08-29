package auth

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/golang-jwt/jwt/v5"
)

func (c impl) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := c.repo.User.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if user.Password != req.Password {
		return nil, model.ErrInvalidCredentials
	}

	now := time.Now()

	// Generate JWT token
	token, err := c.jwtHelper.GenerateJWTToken(map[string]interface{}{
		"sub":  user.ID,
		"iss":  c.cfg.App,
		"role": "user",
		"exp":  jwt.NewNumericDate(now.AddDate(1, 0, 0)),
		"nbf":  jwt.NewNumericDate(now),
		"iat":  jwt.NewNumericDate(now),
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
