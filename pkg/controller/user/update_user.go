package user

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
)

func (c *impl) UpdateUser(ctx context.Context, user model.UpdateUserRequest) (*model.User, error) {
	uID, err := middleware.UserIDFromContext(ctx)
	if err != nil {
		return nil, model.ErrInvalidToken
	}

	u, err := c.repo.User.GetByID(ctx, uID)
	if err != nil {
		return nil, err
	}

	u.Avatar = user.Avatar
	u.FullName = user.FullName

	updated, err := c.repo.User.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return updated, nil
}
