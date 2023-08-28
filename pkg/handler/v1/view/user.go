package view

// MeResponse represent the user response
type MeResponse = Response[Me] // @name MeResponse

// Me represent the user
type Me struct {
	ID    string `json:"id" validate:"required"`
	Email string `json:"email" validate:"required"`
} // @name Me

// UpdateUserRequest represent the update user request
type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
} // @name UpdateUserRequest

// UpdatePasswordRequest represent the update password request
type UpdatePasswordRequest struct {
	Email          string `json:"email" validate:"required"`
	NewPassword    string `json:"new_password" validate:"required"`
	RetypePassword string `json:"retype_password" validate:"required"`
} // @name UpdatePasswordRequest
