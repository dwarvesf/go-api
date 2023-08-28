package viewmodel

// MeResponse represent the user response
type MeResponse struct {
	ID    string `json:"id" validate:"required"`
	Email string `json:"email" validate:"required"`
} // @name MeResponse
