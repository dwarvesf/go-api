package user

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
)

func (c *impl) Me(ctx context.Context) (*model.User, error) {
	uID, err := middleware.UserIDFromContext(ctx)
	if err != nil {
		return nil, model.ErrInvalidToken
	}

	u, err := c.repo.User.GetByID(ctx, uID)
	if err != nil {
		return nil, err
	}

	return u, nil
}
