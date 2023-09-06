package model

// LoginRequest represent the login request
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResponse represent the login response
type LoginResponse struct {
	ID          int
	Email       string
	AccessToken string
}

// SignupRequest represent the signup request
type SignupRequest struct {
	Email          string
	Password       string
	Name           string
	Avatar         string
	Salt           string
	HashedPassword string
	Role           Role
	Status         Status
}
