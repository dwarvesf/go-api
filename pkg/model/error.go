package model

import "net/http"

var (
	// ErrInvalidToken is the error for invalid token
	ErrInvalidToken = Error{
		Status:  http.StatusUnauthorized,
		Code:    "Unauthorized",
		Message: "Unauthorized",
	}

	// ErrUnexpectedAuthorizationHeader is the error for unexpected authorization header
	ErrUnexpectedAuthorizationHeader = Error{
		Status:  http.StatusUnauthorized,
		Code:    "Unauthorized",
		Message: "Unexpected authorization headers",
	}

	// ErrInvalidCredentials is the error for invalid credentials
	ErrInvalidCredentials = Error{
		Status:  http.StatusBadRequest,
		Code:    "WRONG_CREDENTIALS",
		Message: "Wrong username or password",
	}

	// ErrNotFound is the error for not found
	ErrNotFound = Error{
		Status:  http.StatusNotFound,
		Code:    "NOT_FOUND",
		Message: "not found",
	}
)

// Error in server
type Error struct {
	Status  int
	Code    string
	Message string
}

func (e Error) Error() string {
	return e.Message
}

// NewError new a error with message
func NewError(status int, code, msg string) error {
	return Error{
		Status:  status,
		Code:    code,
		Message: msg,
	}
}
