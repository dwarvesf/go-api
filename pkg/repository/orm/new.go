package orm

// Repo represent the repositories
type Repo struct{}

// NewRepo will create an object that represent the Repo interface (repository.Repo)
func NewRepo() Repo {
	return Repo{}
}
