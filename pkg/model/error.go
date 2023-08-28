package model

import "net/http"

var (
	ErrInvalidToken = Error{
		Status:  http.StatusUnauthorized,
		Code:    "Unauthorized",
		Message: "Unauthorized",
	}
	ErrUnexpectedAuthorizationHeader = Error{
		Status:  http.StatusUnauthorized,
		Code:    "Unauthorized",
		Message: "Unexpected authorization headers",
	}
	ErrInvalidCredentials = Error{
		Status:  http.StatusBadRequest,
		Code:    "WRONG_CREDENTIALS",
		Message: "Wrong username or password",
	}
)

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
