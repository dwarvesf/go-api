package realtime

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"sync"

	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/gin-gonic/gin"
)

// SSEHeadersMiddleware sets the headers for SSE.
func SSEHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

// SSEConn represents a SSE connection.
type SSEConn struct {
	Channel chan string
	ID      string
}

type sse struct {
	clients map[string][]*SSEConn
	mutex   sync.Mutex
	authMw  middleware.AuthMiddleware
}

// NewSSE creates a new SSE server.
func NewSSE(authMw middleware.AuthMiddleware) Server {
	return &sse{
		clients: make(map[string][]*SSEConn),
		mutex:   sync.Mutex{},
		authMw:  authMw,
	}
}

func (s *sse) HandleConnection(c *gin.Context) (*User, error) {
	// Create a channel for sending SSE data
	messageChannel := make(chan string)
	var device *SSEConn
	var userID string
	jwtClaims, err := s.authMw.Authenticate(c)
	if err != nil {
		if !errors.Is(err, model.ErrNoAuthHeader) {
			return nil, err
		}
		if errors.Is(err, model.ErrNoAuthHeader) {
			userID = "guest-" + generateRandomID()
			device = &SSEConn{
				Channel: messageChannel,
				ID:      userID,
			}
		}
	} else {
		uID, err := middleware.UserIDFromJWTClaims(jwtClaims)
		if err != nil {
			return nil, err
		}
		userID = "user-" + strconv.Itoa(uID)
		device = &SSEConn{
			Channel: messageChannel,
			ID:      userID + "-" + generateRandomID(),
		}
	}

	// Register the client's channel for SSE updates
	s.mutex.Lock()
	clientArr, ok := s.clients[userID]
	if !ok {
		clientArr = make([]*SSEConn, 0)
	}
	s.clients[userID] = append(clientArr, device)
	s.mutex.Unlock()

	user := &User{
		ID:       userID,
		DeviceID: device.ID,
	}

	return user, nil
}

func (s *sse) HandleEvent(c *gin.Context, u User, fn func(any) error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientChArr, ok := s.clients[u.ID]
	if !ok {
		return
	}

	var clientCh *SSEConn
	for _, ch := range clientChArr {
		if ch.ID == u.DeviceID {
			clientCh = ch
		}
	}

	go func() {
		disconnected := c.Stream(func(w io.Writer) bool {
			// Stream message to client from message channel
			if msg, ok := <-clientCh.Channel; ok {
				c.SSEvent("message", msg)
				return true
			}
			return false
		})

		if disconnected {
			s.DisconnectUser(u)
		}
	}()
}
func (s *sse) SendMessage(userID string, message string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientChArr, ok := s.clients[userID]
	if !ok {
		return errors.New("client not found")
	}

	for _, clientCh := range clientChArr {
		clientCh.Channel <- message
	}

	return nil
}

func (s *sse) SendData(userID string, data any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	clientChArr, ok := s.clients[userID]
	if !ok {
		return errors.New("client not found")
	}

	for _, clientCh := range clientChArr {
		clientCh.Channel <- string(body)
	}

	return nil
}

func (s *sse) BroadcastMessage(message string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, chArr := range s.clients {
		for _, clientCh := range chArr {
			clientCh.Channel <- message
		}
	}

	return nil
}

func (s *sse) BroadcastData(data any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for _, chArr := range s.clients {
		for _, clientCh := range chArr {
			clientCh.Channel <- string(body)
		}
	}

	return nil
}

func (s *sse) DisconnectUser(u User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientChArr, ok := s.clients[u.ID]
	if !ok {
		return errors.New("client not found")
	}

	for _, clientCh := range clientChArr {
		close(clientCh.Channel)
	}

	delete(s.clients, u.ID)

	return nil
}
