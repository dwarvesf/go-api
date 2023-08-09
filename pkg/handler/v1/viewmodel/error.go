package viewmodel

import (
	"net/http"

	"github.com/dwarvesf/go-api/pkg/model"
)

// ErrorResponse in server
type ErrorResponse struct {
	Status  int    `json:"status" validate:"required"`
	Code    string `json:"code" validate:"required"`
	Message string `json:"message" validate:"required"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

// NewError new a error with message
func NewError(status int, code, msg string) error {
	return ErrorResponse{
		Status:  status,
		Code:    code,
		Message: msg,
	}
}

// ErrBadRequest new a bad request error with the detail
func ErrBadRequest(err error) error {
	return model.Error{
		Status:  http.StatusBadRequest,
		Code:    "BAD_REQUEST",
		Message: err.Error(),
	}
}
