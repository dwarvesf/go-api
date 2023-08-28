package viewmodel

// LoginRequest represent the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
} // @name LoginRequest

// LoginResponse represent the login response
type LoginResponse struct {
	ID          string `json:"id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	AccessToken string `json:"access_token" validate:"required"`
} // @name LoginResponse

type SignupRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
} // @name SignupRequest
