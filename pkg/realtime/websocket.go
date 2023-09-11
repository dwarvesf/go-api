package realtime

import (
	"strconv"
	"sync"

	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Socket represents a WebSocket connection
type Socket interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	WriteJSON(v interface{}) error
	Close() error
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Conn represents a WebSocket connection
type Conn struct {
	Socket
	DeviceID    string
	IsGuest     bool
	Permissions []string
}

type ws struct {
	clients map[string][]*Conn
	mutex   sync.RWMutex
	authMw  middleware.AuthMiddleware
	log     logger.Log
}

// HandleConnection handles WebSocket connections and user authentication.
func (s *ws) HandleConnection(c *gin.Context) (*User, error) {
	r := c.Request
	w := c.Writer
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	var device *Conn
	var userID string
	jwtClaims, err := s.authMw.Authenticate(c)
	if err != nil {
		if !errors.Is(err, model.ErrNoAuthHeader) {
			return nil, err
		}

		if errors.Is(err, model.ErrNoAuthHeader) {
			userID = PrefixGuest + generateRandomID()
			device = &Conn{
				DeviceID: userID,
				IsGuest:  true,
				Socket:   conn,
			}
		}
	} else {
		uID, err := middleware.UserIDFromJWTClaims(jwtClaims)
		if err != nil {
			return nil, err
		}
		userID = PrefixUser + strconv.Itoa(uID)
		device = &Conn{
			DeviceID: userID + "-" + generateRandomID(),
			Socket:   conn,
		}
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Add the user to the server's list of clients
	if _, found := s.clients[userID]; !found {
		s.clients[userID] = make([]*Conn, 0)
	}
	s.clients[userID] = append(s.clients[userID], device)

	return &User{
		ID:       userID,
		DeviceID: device.DeviceID,
	}, nil
}

func (s *ws) HandleEvent(c *gin.Context, u User, callback func(*gin.Context, any) error) {
	var conn *Conn
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	devices, ok := s.clients[u.ID]
	if !ok {
		return
	}
	for _, device := range devices {
		if device.DeviceID == u.DeviceID {
			conn = device
			break
		}
	}

	defer s.DisconnectUser(u)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			s.log.Error(err)
			return
		}

		err = callback(c, message)
		if err != nil {
			s.log.Error(err)
			return
		}
	}
}

// SendMessage sends a message to all devices of a WebSocket user.
func (s *ws) SendMessage(userID string, message string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	devices, found := s.clients[userID]
	if !found {
		return ErrUserNotFound
	}

	for _, device := range devices {
		err := device.Socket.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			return err
		}
	}

	return nil
}

// SendData sends data to all devices of a WebSocket user.
func (s *ws) SendData(userID string, data any) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	devices, found := s.clients[userID]
	if !found {
		return ErrUserNotFound
	}

	for _, device := range devices {
		err := device.Socket.WriteJSON(data)
		if err != nil {
			return err
		}
	}

	return nil
}

// BroadcastMessage sends a message to all devices of all WebSocket users.
func (s *ws) BroadcastMessage(message string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, devices := range s.clients {
		for _, device := range devices {
			err := device.Socket.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// BroadcastData sends data to all devices of all WebSocket users.
func (s *ws) BroadcastData(data any) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, devices := range s.clients {
		for _, device := range devices {
			err := device.Socket.WriteJSON(data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ws) DisconnectUser(u User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	devices, found := s.clients[u.ID]
	if !found {
		return ErrUserNotFound
	}

	for i, device := range devices {
		if device.DeviceID == u.DeviceID {
			device.Close()
			devices = append(devices[:i], devices[i+1:]...)
			return nil
		}
	}

	return ErrDeviceNotFound
}
