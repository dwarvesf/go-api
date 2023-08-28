package model

// UpdateUserRequest represent the update user request
type UpdateUserRequest struct {
	FullName string
	Status   string
	Avatar   string
}

// UpdatePasswordRequest represent the update password request
type UpdatePasswordRequest struct {
	Email          string
	NewPassword    string
	RetypePassword string
}
