package monitor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock "github.com/dwarvesf/go-api/mocks/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSentryPanicMiddleware(t *testing.T) {
	testCases := map[string]struct {
		panicMessage     string
		expectedTraceID  string
		expectedHTTPCode int
	}{
		"success": {
			panicMessage:     "test panic",
			expectedTraceID:  "00000000000000000000000000000000",
			expectedHTTPCode: http.StatusInternalServerError,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := gin.New()
			r.Use(SentryPanicMiddleware(mock.NewLog(t)))

			// Create a test request and recorder
			req, _ := http.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Perform a request that triggers a panic
			r.GET("/test", func(c *gin.Context) {
				panic(tc.panicMessage)
			})

			r.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedHTTPCode, w.Code)
			assert.JSONEq(t, `{
				"status": 500,
				"code": "INTERNAL_SERVER_ERROR",
				"message": "Internal Server Error",
				"traceID": "`+tc.expectedTraceID+`"
			}`, w.Body.String())
		})
	}
}
