package viewmodel

// MessageResponse is the response for message
type MessageResponse struct {
	Message string `json:"message" validate:"required"`
} // @name MessageResponse
