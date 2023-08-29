package user

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/model"
)

// Repo represent the user
type Repo interface {
	GetByID(ctx context.Context, id int) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.SignupRequest) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	UpdatePassword(ctx context.Context, uID int, newPassword string) error
}

// New return new user repo
func New() Repo {
	return &mem{
		users: make(map[int]model.User),
	}
}
