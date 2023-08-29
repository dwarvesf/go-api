package user

import (
	"context"

	"github.com/dwarvesf/go-api/pkg/model"
)

type mem struct {
	users map[int]model.User
}

func (m mem) GetByID(ctx context.Context, id int) (*model.User, error) {
	if user, ok := m.users[id]; ok {
		return &user, nil
	}

	return nil, model.ErrNotFound
}

func (m mem) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, model.ErrNotFound
}

func (m *mem) Create(ctx context.Context, user *model.SignupRequest) (*model.User, error) {
	id := len(m.users) + 1
	u := model.User{
		ID:       id,
		Email:    user.Email,
		Password: user.Password,
		FullName: user.FullName,
		Status:   "active",
		Avatar:   "",
	}

	m.users[id] = u

	return &u, nil
}

func (m *mem) Update(ctx context.Context, user *model.User) (*model.User, error) {
	m.users[user.ID] = *user
	return user, nil
}

func (m *mem) UpdatePassword(ctx context.Context, uID int, newPassword string) error {
	if u, ok := m.users[uID]; ok {
		u.Password = newPassword
		m.users[uID] = u
		return nil
	}

	return model.ErrNotFound
}
