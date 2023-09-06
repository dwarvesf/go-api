package user

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/db"
)

func (c impl) UpdatePassword(ctx context.Context, user model.UpdatePasswordRequest) error {
	uID, err := middleware.UserIDFromContext(ctx)
	if err != nil {
		return model.ErrInvalidToken
	}
	dbCtx := db.FromContext(ctx)
	u, err := c.repo.User.GetByID(dbCtx, uID)
	if err != nil {
		return err
	}

	if u.HashedPassword != user.OldPassword {
		return model.ErrInvalidCredentials
	}

	err = c.repo.User.UpdatePassword(dbCtx, uID, user.NewPassword)
	return err
}
