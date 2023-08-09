package util

import (
	"errors"
	"net/http"

	"github.com/dwarvesf/go-api/pkg/handler/v1/viewmodel"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/gin-gonic/gin"
)

// HandleError handle the rest error
func HandleError(c *gin.Context, err error) {
	e := tryParseError(err)
	c.JSON(e.Status, gin.H{
		"status":  e.Status,
		"code":    e.Code,
		"message": e.Message,
	})

}

func tryParseError(err error) viewmodel.ErrorResponse {
	var e model.Error
	ok := errors.As(err, &e)
	if ok {
		return viewmodel.ErrorResponse{
			Status:  e.Status,
			Code:    e.Code,
			Message: e.Message,
		}
	}

	var viewErr viewmodel.ErrorResponse
	ok = errors.As(err, &viewErr)
	if ok {
		return viewErr
	}

	return viewmodel.ErrorResponse{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_ERROR",
		Message: err.Error(),
	}
}
