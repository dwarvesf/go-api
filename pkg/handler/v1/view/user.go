package view

// MeResponse represent the user response
type MeResponse = Response[Me] // @name MeResponse

// Me represent the user
type Me struct {
	ID    int    `json:"id" validate:"required"`
	Email string `json:"email" validate:"required"`
} // @name Me

// UpdateUserRequest represent the update user request
type UpdateUserRequest struct {
	FullName string `json:"fullName"`
	Avatar   string `json:"avatar"`
} // @name UpdateUserRequest

// UpdatePasswordRequest represent the update password request
type UpdatePasswordRequest struct {
	NewPassword string `json:"newPassword" validate:"required"`
	OldPassword string `json:"oldPassword" validate:"required"`
} // @name UpdatePasswordRequest

// UserResponse represent the user response
type UserResponse = Response[User] // @name UserResponse

// User represent the user
type User struct {
	ID       int    `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Avatar   string `json:"avatar" validate:"required"`
} // @name User
