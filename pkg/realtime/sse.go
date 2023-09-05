package realtime

import (
	"encoding/json"
	"errors"
	"io"
	"sync"

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

type sse struct {
	clients map[string][]chan<- string
	mutex   sync.Mutex
}

// NewSSE creates a new SSE server.
func NewSSE() Server {
	return &sse{}
}

func (s *sse) HandleConnection(c *gin.Context) error {
	r := c.Request

	// Create a channel for sending SSE data
	messageChannel := make(chan string)

	// Register the client's channel for SSE updates
	clientID := r.RemoteAddr // You can use a unique client identifier
	s.mutex.Lock()
	clientArr, ok := s.clients[clientID]
	if !ok {
		clientArr = make([]chan<- string, 0)
	}
	s.clients[clientID] = append(clientArr, messageChannel)
	s.mutex.Unlock()

	c.Stream(func(w io.Writer) bool {
		// Stream message to client from message channel
		if msg, ok := <-messageChannel; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})

	return nil
}
func (s *sse) SendMessage(userID string, message string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientChArr, ok := s.clients[userID]
	if !ok {
		return errors.New("client not found")
	}
	for _, ch := range clientChArr {
		ch <- message
	}

	return nil

}
func (s *sse) SendData(userID string, data any) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	clientChArr, ok := s.clients[userID]
	if !ok {
		return errors.New("client not found")
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for _, clientCh := range clientChArr {
		clientCh <- string(body)
	}

	return nil
}
func (s *sse) BroadcastMessage(message string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, chArr := range s.clients {
		for _, ch := range chArr {
			ch <- message
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
			clientCh <- string(body)
		}
	}

	return nil
}
