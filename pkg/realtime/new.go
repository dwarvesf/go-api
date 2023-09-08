package realtime

import (
	"sync"

	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/util"
	"github.com/gin-gonic/gin"
)

// User represents a realtime user
type User struct {
	ID       string
	DeviceID string
}

// Server represents a WebSocket server interface
type Server interface {
	HandleConnection(c *gin.Context) (*User, error)
	HandleEvent(c *gin.Context, u User, callback func(data any) error)
	SendMessage(userID string, message string) error
	SendData(userID string, data any) error
	BroadcastMessage(message string) error
	BroadcastData(data any) error
	DisconnectUser(u User) error
}

// generateRandomID generates a random ID for guest users.
func generateRandomID() string {
	return util.RandomString(8)
}

// New creates a new WebSocket server.
func New(authMw middleware.AuthMiddleware) Server {
	return &impl{
		Clients: make(map[string][]*Conn),
		Mutex:   sync.RWMutex{},
		authMw:  authMw,
	}
}
