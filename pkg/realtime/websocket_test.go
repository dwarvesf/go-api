package realtime

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSocket struct {
	content []byte
}

func (m mockSocket) ReadMessage() (messageType int, p []byte, err error) {
	return 0, m.content, nil
}

func (m *mockSocket) WriteMessage(messageType int, data []byte) error {
	m.content = append(m.content, data...)
	return nil
}

func (m *mockSocket) WriteJSON(v interface{}) error {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	m.content = append(m.content, jsonBytes...)
	return nil
}

func (m *mockSocket) Close() error {
	m.content = []byte{}
	return nil
}

func (m *mockSocket) Clear() error {
	m.content = []byte{}
	return nil
}

func Test_ws_BroadcastMessage(t *testing.T) {
	message := "Hello, World!"

	server := &ws{
		clients: make(map[string][]*Conn),
		mutex:   sync.RWMutex{},
	}

	userID := "testUser"
	mockConnection := &mockSocket{}

	conn := &Conn{
		Socket:   mockConnection,
		DeviceID: "device1",
	}
	server.clients[userID] = []*Conn{conn}

	err := server.BroadcastMessage(message)

	assert.Nil(t, err)
	assert.Equal(t, message, string(mockConnection.content))
}

func Test_ws_SendMessage(t *testing.T) {
	s := &ws{
		clients: make(map[string][]*Conn),
		mutex:   sync.RWMutex{},
	}

	// Create a dummy client
	mockConnection := &mockSocket{}
	s.clients["user1"] = []*Conn{
		{
			Socket:   mockConnection,
			DeviceID: "user1-client1",
		},
	}

	tests := map[string]struct {
		userID  string
		message string
		wantErr bool
	}{
		"success": {
			userID:  "user1",
			message: "Hello, World!",
			wantErr: false,
		},
		"failure": {
			userID:  "nonExistentUser",
			message: "",
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockConnection.Clear()
			err := s.SendMessage(tc.userID, tc.message)

			if tc.wantErr == (err == nil) {
				t.Errorf("%v case: expected error %v but got %v", name, tc.wantErr, err)
				return
			}

			if string(mockConnection.content) != tc.message {
				t.Errorf("%v case: expected message %q but got %q", name, tc.message, string(mockConnection.content))
			}
		})
	}

}

func Test_ws_SendData(t *testing.T) {
	type testStruct struct {
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
	s := &ws{
		clients: make(map[string][]*Conn),
		mutex:   sync.RWMutex{},
	}

	// Create a dummy client
	mockConnection := &mockSocket{}
	s.clients["user1"] = []*Conn{
		{
			Socket:   mockConnection,
			DeviceID: "user1-client1",
		},
	}

	tests := map[string]struct {
		userID  string
		data    any
		wantErr bool
		message string
	}{
		"success": {
			userID:  "user1",
			data:    testStruct{Name: "user1", Age: 20},
			message: `{"name":"user1","age":20}`,
			wantErr: false,
		},
		"client not found": {
			userID:  "user2",
			data:    testStruct{Name: "user2", Age: 20},
			wantErr: true,
		},
		"invalid data": {
			userID:  "user1",
			data:    make(chan string),
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockConnection.Clear()
			err := s.SendData(tc.userID, tc.data)
			if (err != nil) != tc.wantErr {
				t.Fatalf("%v case: SendData() error = %v, wantErr %v", name, err, tc.wantErr)
				return
			}

			if string(mockConnection.content) != tc.message {
				t.Errorf("%v case: expected message %q but got %q", name, tc.message, string(mockConnection.content))
			}
		})
	}
}

func Test_ws_BroadcastData(t *testing.T) {
	type testStruct struct {
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
	s := &ws{
		clients: make(map[string][]*Conn),
		mutex:   sync.RWMutex{},
	}

	// Create a dummy client
	mockConnection := &mockSocket{}
	s.clients["user1"] = []*Conn{
		{
			Socket:   mockConnection,
			DeviceID: "user1-client1",
		},
	}

	tests := map[string]struct {
		data    any
		wantErr bool
		message string
	}{
		"success": {
			data:    testStruct{Name: "user1", Age: 20},
			message: `{"name":"user1","age":20}`,
			wantErr: false,
		},
		"failure": {
			data:    make(chan string),
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockConnection.Clear()
			err := s.BroadcastData(tc.data)
			if (err != nil) != tc.wantErr {
				t.Fatalf("%v case: BroadcastData() error = %v, wantErr %v", name, err, tc.wantErr)
			}
			if string(mockConnection.content) != tc.message {
				t.Errorf("%v case: expected message %q but got %q", name, tc.message, string(mockConnection.content))
			}
		})
	}
}
