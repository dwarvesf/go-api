package view

// MessageResponse is the response for message
type MessageResponse = Response[Message] // @name MessageResponse

// Message is the message
type Message struct {
	Message string `json:"message" validate:"required"`
} // @name Message

// Response is the response for data
type Response[T any] struct {
	Data T `json:"data"`
} // @name Response

// ListResponse is the response for list data
type ListResponse[T any] struct {
	Data     []T      `json:"data"`
	Metadata Metadata `json:"metadata"`
} // @name ListResponse
