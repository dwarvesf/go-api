package model

// LoginRequest represent the login request
type LoginRequest struct {
	Email    string
	Password string
}

// LoginResponse represent the login response
type LoginResponse struct {
	ID          string
	Email       string
	AccessToken string
}

// SignupRequest represent the signup request
type SignupRequest struct {
	Email    string
	Password string
	FullName string
	Status   string
	Avatar   string
}
