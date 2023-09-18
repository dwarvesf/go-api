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
	clients sync.Map
	authMw  middleware.AuthMiddleware
}

// NewSSE creates a new SSE server.
func NewSSE(authMw middleware.AuthMiddleware) Server {
	return &sse{
		clients: sync.Map{},
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
			userID = PrefixGuest + generateRandomID()
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
		userID = PrefixUser + strconv.Itoa(uID)
		device = &SSEConn{
			Channel: messageChannel,
			ID:      userID + "-" + generateRandomID(),
		}
	}

	// Register the client's channel for SSE updates
	clientArr, ok := s.clients.Load(userID)
	if !ok {
		clientArr = make(map[string]*SSEConn, 0)
	}
	clientArrData, ok := clientArr.(map[string]*SSEConn)
	if !ok {
		clientArrData = make(map[string]*SSEConn, 0)
	}
	clientArrData[device.ID] = device
	s.clients.Store(userID, clientArrData)

	user := &User{
		ID:       userID,
		DeviceID: device.ID,
	}

	return user, nil
}

func (s *sse) HandleEvent(c *gin.Context, u User, callback func(*gin.Context, any) error) {
	val, ok := s.clients.Load(u.ID)
	if !ok {
		return
	}
	clientChArr, ok := val.(map[string]*SSEConn)
	if !ok {
		return
	}

	var clientCh *SSEConn
	for _, ch := range clientChArr {
		if ch.ID == u.DeviceID {
			clientCh = ch
			break
		}
	}

	finished := make(chan bool)

	go func() {
		close(finished)
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
	<-finished
}
func (s *sse) SendMessage(userID string, message string) error {
	val, ok := s.clients.Load(userID)
	if !ok {
		return ErrClientNotFound
	}
	clientChArr, ok := val.(map[string]*SSEConn)
	if !ok {
		return ErrClientNotFound
	}

	for _, clientCh := range clientChArr {
		clientCh.Channel <- message
	}

	return nil
}

func (s *sse) SendData(userID string, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	val, ok := s.clients.Load(userID)
	if !ok {
		return ErrClientNotFound
	}
	clientChArr, ok := val.(map[string]*SSEConn)
	if !ok {
		return ErrClientNotFound
	}

	for _, clientCh := range clientChArr {
		clientCh.Channel <- string(body)
	}

	return nil
}

func (s *sse) BroadcastMessage(message string) error {
	s.clients.Range(func(key, value any) bool {
		clientChArr, ok := value.(map[string]*SSEConn)
		if !ok {
			return true
		}
		for i := range clientChArr {
			go func(key string) {
				clientChArr[key].Channel <- message
			}(i)
		}
		return true
	})

	return nil
}

func (s *sse) BroadcastData(data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	s.clients.Range(func(key, value any) bool {
		clientChArr, ok := value.(map[string]*SSEConn)
		if !ok {
			return true
		}
		for _, clientCh := range clientChArr {
			clientCh.Channel <- string(body)
		}
		return true
	})

	return nil
}

func (s *sse) DisconnectUser(u User) error {
	val, ok := s.clients.Load(u.ID)
	if !ok {
		return nil
	}
	clientChArr, ok := val.(map[string]*SSEConn)
	if !ok {
		return nil
	}

	clientCh, ok := clientChArr[u.DeviceID]
	if !ok {
		return nil
	}

	close(clientCh.Channel)
	delete(clientChArr, u.DeviceID)

	s.clients.Store(u.ID, clientChArr)
	return nil
}
