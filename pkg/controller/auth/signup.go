package auth

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/model"
)

func (c impl) Signup(ctx context.Context, req model.SignupRequest) error {
	_, err := c.repo.User.Create(ctx, &req)
	if err != nil {
		return err
	}

	return nil
}
