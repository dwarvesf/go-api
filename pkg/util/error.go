package util

import (
	"errors"
	"net/http"

	"github.com/dwarvesf/go-api/pkg/handler/v1/view"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/monitor"
	"github.com/gin-gonic/gin"
)

// HandleError handle the rest error
func HandleError(c *gin.Context, err error) {
	e := tryParseError(err)
	c.JSON(e.Status, gin.H{
		"status":  e.Status,
		"code":    e.Code,
		"message": e.Err,
		"traceID": monitor.GetTraceID(c.Request.Context()),
	})
}

func tryParseError(err error) view.ErrorResponse {
	var e model.Error
	ok := errors.As(err, &e)
	if ok {
		return view.ErrorResponse{
			Status: e.Status,
			Code:   e.Code,
			Err:    e.Message,
		}
	}

	var viewErr view.ErrorResponse
	ok = errors.As(err, &viewErr)
	if ok {
		return viewErr
	}

	return view.ErrorResponse{
		Status: http.StatusInternalServerError,
		Code:   "INTERNAL_ERROR",
		Err:    err.Error(),
	}
}
