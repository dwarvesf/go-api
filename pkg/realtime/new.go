package realtime

import (
	"github.com/dwarvesf/go-api/pkg/util"
	"github.com/gin-gonic/gin"
)

// Server represents a WebSocket server interface
type Server interface {
	HandleConnection(c *gin.Context) error
	SendMessage(userID string, message string) error
	SendData(userID string, data any) error
	BroadcastMessage(message string) error
	BroadcastData(data any) error
}

// generateRandomID generates a random ID for guest users.
func generateRandomID() string {
	return util.RandomString(8)
}
