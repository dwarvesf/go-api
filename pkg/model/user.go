package model

// UpdateUserRequest represent the update user request
type UpdateUserRequest struct {
	FullName string
	Avatar   string
}

// UpdatePasswordRequest represent the update password request
type UpdatePasswordRequest struct {
	OldPassword string
	NewPassword string
}

// User represent the user
type User struct {
	ID       int
	Email    string
	Password string
	FullName string
	Status   string
	Avatar   string
}
