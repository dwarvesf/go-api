package repository

import "github.com/dwarvesf/go-api/pkg/repository/user"

// Repo represent the repository
type Repo struct {
	User user.Repo
}

// NewRepo will create an object that represent the Repo interface
func NewRepo() *Repo {
	return &Repo{
		User: user.New(),
	}
}
