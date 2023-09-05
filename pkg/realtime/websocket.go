package realtime

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Conn represents a WebSocket connection
type Conn struct {
	*websocket.Conn
	ID          string
	IsGuest     bool
	Permissions []string
}

type impl struct {
	Clients map[string][]*Conn
	Mutex   sync.Mutex
}

// New creates a new WebSocket server.
func New() Server {
	return &impl{
		Clients: make(map[string][]*Conn),
	}
}

// HandleConnection handles WebSocket connections and user authentication.
func (s *impl) HandleConnection(c *gin.Context) error {
	r := c.Request
	w := c.Writer
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	// Extract the user's token
	token := r.Header.Get("Authorization")

	var user *Conn
	if token == "" {
		// Guest user
		user = &Conn{
			ID:      "guest-" + generateRandomID(), // Generate a unique ID for guest
			IsGuest: true,
			Conn:    conn,
		}
	} else {
		// Perform token validation and role determination (user or admin)
		// Example: Use a middleware to validate and set user roles and permissions
		permissions := []string{} // Determine user permissions based on token

		user = &Conn{
			ID:          "user_id_from_token",
			Permissions: permissions,
			Conn:        conn,
		}
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// Add the user to the server's list of clients
	s.Clients[user.ID] = append(s.Clients[user.ID], user)

	return nil
}

// SendMessage sends a message to all devices of a WebSocket user.
func (s *impl) SendMessage(userID string, message string) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

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
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

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
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

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
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

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
