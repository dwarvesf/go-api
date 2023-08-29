package view

// LoginRequest represent the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
} // @name LoginRequest

// LoginResponse represent the login response
type LoginResponse = Response[Auth] // @name LoginResponse

// Auth represent the auth response
type Auth struct {
	ID          int    `json:"id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	AccessToken string `json:"accessToken" validate:"required" `
} // @name Auth

// SignupRequest represent the signup request
type SignupRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"fullName" binding:"required"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
} // @name SignupRequest
