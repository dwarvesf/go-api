package realtime

import (
	"sync"

	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/util"
	"github.com/gin-gonic/gin"
)

const (
	// PrefixUser is the prefix for authenticated users
	PrefixUser = "user-"

	// PrefixGuest is the prefix for guest users
	PrefixGuest = "guest-"

	// randomIDLength is the length of the random ID for guest users
	randomIDLength = 10
)

// User represents a realtime user
type User struct {
	ID       string
	DeviceID string
}

// Server represents a WebSocket server interface
type Server interface {
	HandleConnection(c *gin.Context) (*User, error)
	HandleEvent(c *gin.Context, u User, callback func(c *gin.Context, data any) error)
	SendMessage(userID string, message string) error
	SendData(userID string, data any) error
	BroadcastMessage(message string) error
	BroadcastData(data any) error
	DisconnectUser(u User) error
}

// generateRandomID generates a random ID for guest users.
func generateRandomID() string {
	return util.RandomString(randomIDLength)
}

// New creates a new WebSocket server.
func New(authMw middleware.AuthMiddleware, l logger.Log) Server {
	return &ws{
		clients: make(map[string]map[string]*Conn),
		mutex:   sync.RWMutex{},
		authMw:  authMw,
		log:     l,
	}
}
