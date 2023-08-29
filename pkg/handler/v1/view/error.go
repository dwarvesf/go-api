package view

import (
	"net/http"

	"github.com/dwarvesf/go-api/pkg/model"
)

// ErrorResponse error response
type ErrorResponse struct {
	Status  int           `json:"-"`
	Err     string        `json:"error" validate:"required"`
	Code    string        `json:"code" validate:"required"`
	TraceID string        `json:"traceId" validate:"required"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
} // @name ErrorResponse

// ErrorDetail error detail
type ErrorDetail struct {
	Field string `json:"field" validate:"required"`
	Error string `json:"error" validate:"required"`
} // @name ErrorDetail

func (e ErrorResponse) Error() string {
	return e.Err
}

// NewError new a error with message
func NewError(status int, code, msg string) error {
	return ErrorResponse{
		Status: status,
		Code:   code,
		Err:    msg,
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
