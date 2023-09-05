package user

import (
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
)

// Repo represent the user
type Repo interface {
	GetList(ctx db.Context, page int, pageSize int, sort string, query string) (*model.UserList, error)
	Count(ctx db.Context) (int64, error)
	GetByID(ctx db.Context, id int) (*model.User, error)
	GetByEmail(ctx db.Context, email string) (*model.User, error)
	Create(ctx db.Context, user model.SignupRequest) (*model.User, error)
	Update(ctx db.Context, uID int, user model.UpdateUserRequest) (*model.User, error)
	UpdatePassword(ctx db.Context, uID int, newPassword string) error
}

// New return new user repo
func New() Repo {
	return &repo{}
}

func toUserModel(user *orm.User) *model.User {
	if user == nil {
		return nil
	}
	return &model.User{
		ID:             user.ID,
		Email:          user.Email,
		FullName:       user.Name,
		Status:         user.Status,
		Avatar:         user.Avatar,
		HashedPassword: user.HashedPassword,
		Role:           user.Role,
		Salt:           user.Salt,
	}
}
