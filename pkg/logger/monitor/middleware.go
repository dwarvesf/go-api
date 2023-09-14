package monitor

import (
	"net/http"
	"time"

	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// SentryPanicMiddleware return the middleware
func SentryPanicMiddleware(log logger.Log) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				sentry.CurrentHub().Recover(err)
				sentry.Flush(time.Second * 5)

				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"code":    "INTERNAL_SERVER_ERROR",
					"message": "Internal Server Error",
					"traceID": GetTraceID(c.Request.Context()),
				})
			}
		}()

		c.Next()
	}
}
