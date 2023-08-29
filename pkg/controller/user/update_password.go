package user

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
)

func (c impl) UpdatePassword(ctx context.Context, user model.UpdatePasswordRequest) error {
	uID, err := middleware.UserIDFromContext(ctx)
	if err != nil {
		return model.ErrInvalidToken
	}
	u, err := c.repo.User.GetByID(ctx, uID)
	if err != nil {
		return err
	}

	if u.Password != user.OldPassword {
		return model.ErrInvalidCredentials
	}

	err = c.repo.User.UpdatePassword(ctx, uID, user.NewPassword)
	return err
}
