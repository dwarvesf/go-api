package realtime

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Conn represents a WebSocket connection
type Conn struct {
	*websocket.Conn
	DeviceID    string
	IsGuest     bool
	Permissions []string
}

type impl struct {
	Clients map[string][]*Conn
	Mutex   sync.RWMutex
	authMw  middleware.AuthMiddleware
}

// HandleConnection handles WebSocket connections and user authentication.
func (s *impl) HandleConnection(c *gin.Context) (*User, error) {
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
			userID = "guest-" + generateRandomID()
			device = &Conn{
				DeviceID: userID,
				IsGuest:  true,
				Conn:     conn,
			}
		}
	} else {
		uID, err := middleware.UserIDFromJWTClaims(jwtClaims)
		if err != nil {
			return nil, err
		}
		userID = "user-" + strconv.Itoa(uID)
		device = &Conn{
			DeviceID: userID + "-" + generateRandomID(),
			Conn:     conn,
		}
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Add the user to the server's list of clients
	if _, found := s.Clients[userID]; !found {
		s.Clients[userID] = make([]*Conn, 0)
	}
	s.Clients[userID] = append(s.Clients[userID], device)

	return &User{
		ID:       userID,
		DeviceID: device.DeviceID,
	}, nil
}

func (s *impl) HandleEvent(c *gin.Context, u User, callback func(data any) error) {
	var conn *Conn
	s.Mutex.RLock()
	devices, ok := s.Clients[u.ID]
	if !ok {
		s.Mutex.RUnlock()
		return
	}
	for _, device := range devices {
		if device.DeviceID == u.DeviceID {
			conn = device
			break
		}
	}
	s.Mutex.RUnlock()

	defer s.DisconnectUser(u)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = callback(message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// SendMessage sends a message to all devices of a WebSocket user.
func (s *impl) SendMessage(userID string, message string) error {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	devices, found := s.Clients[userID]
	if !found {
		return fmt.Errorf("user not found")
	}

	for _, device := range devices {
		err := device.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			return err
		}
	}

	return nil
}

// SendData sends data to all devices of a WebSocket user.
func (s *impl) SendData(userID string, data any) error {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	devices, found := s.Clients[userID]
	if !found {
		return fmt.Errorf("user not found")
	}

	for _, device := range devices {
		err := device.Conn.WriteJSON(data)
		if err != nil {
			return err
		}
	}

	return nil
}

// BroadcastMessage sends a message to all devices of all WebSocket users.
func (s *impl) BroadcastMessage(message string) error {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	for _, devices := range s.Clients {
		for _, device := range devices {
			err := device.Conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

// BroadcastData sends data to all devices of all WebSocket users.
func (s *impl) BroadcastData(data any) error {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	for _, devices := range s.Clients {
		for _, device := range devices {
			err := device.Conn.WriteJSON(data)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

func (s *impl) DisconnectUser(u User) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	devices, found := s.Clients[u.ID]
	if !found {
		return fmt.Errorf("user not found")
	}

	for i, device := range devices {
		if device.DeviceID == u.DeviceID {
			device.Close()
			devices = append(devices[:i], devices[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("device not found")
}
