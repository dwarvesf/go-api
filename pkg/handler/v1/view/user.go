package view

// MeResponse represent the user response
type MeResponse = Response[Me] // @name MeResponse

// Me represent the user
type Me struct {
	ID    string `json:"id" validate:"required"`
	Email string `json:"email" validate:"required"`
} // @name Me
