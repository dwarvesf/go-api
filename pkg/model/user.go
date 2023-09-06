package model

// Role represent the user role
type Role string

const (
	// RoleUser is the user role
	RoleUser Role = "user"
	// RoleAdmin is the admin role
	RoleAdmin Role = "admin"
)

// Status represent the user status
type Status string

const (
	// StatusActive is the active status
	StatusActive Status = "active"
	// StatusInactive is the inactive status
	StatusInactive Status = "inactive"
)

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
	ID             int
	Email          string
	HashedPassword string
	Salt           string
	FullName       string
	Status         string
	Avatar         string
	Role           string
}

// UserList represent the user list
type UserList struct {
	Pagination Pagination
	Data       []*User
}
