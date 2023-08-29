package user

import (
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
)

type mem struct {
	users map[int]model.User
}

func (m mem) GetByID(ctx db.Context, uID int) (*model.User, error) {
	orm.FindUser(ctx.Context, ctx.DB, uID)
	if user, ok := m.users[uID]; ok {
		return &user, nil
	}

	return nil, model.ErrNotFound
}

func (m mem) GetByEmail(ctx db.Context, email string) (*model.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, model.ErrNotFound
}

func (m *mem) Create(ctx db.Context, user *model.SignupRequest) (*model.User, error) {
	id := len(m.users) + 1
	u := model.User{
		ID:             id,
		Email:          user.Email,
		HashedPassword: user.Password,
		FullName:       user.Name,
		Status:         "active",
		Avatar:         "",
	}

	m.users[id] = u

	return &u, nil
}

func (m *mem) Update(ctx db.Context, user *model.User) (*model.User, error) {
	m.users[user.ID] = *user
	return user, nil
}

func (m *mem) UpdatePassword(ctx db.Context, uID int, newPassword string) error {
	if u, ok := m.users[uID]; ok {
		u.HashedPassword = newPassword
		m.users[uID] = u
		return nil
	}

	return model.ErrNotFound
}
