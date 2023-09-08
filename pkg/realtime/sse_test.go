package realtime

import (
	"net/http/httptest"
	"sync"
	"testing"
)

type TestResponseRecorder struct {
	*httptest.ResponseRecorder
	closeChannel chan bool
}

func (r *TestResponseRecorder) CloseNotify() <-chan bool {
	return r.closeChannel
}

func (r *TestResponseRecorder) closeClient() {
	r.closeChannel <- true
}

// func TestHandleConnection(t *testing.T) {
// 	tests := map[string]struct {
// 		clientID string
// 		message  string
// 		wantErr  bool
// 	}{
// 		"success": {
// 			clientID: "client1",
// 			message:  "event:message\ndata:test message\n\n",
// 			wantErr:  false,
// 		},
// 		"failure": {
// 			clientID: "client2",
// 			wantErr:  true,
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			closeChannel := make(chan bool)
// 			w := &TestResponseRecorder{
// 				httptest.NewRecorder(),
// 				closeChannel,
// 			}
// 			ginCtx := testutil.NewRequest(w, testutil.MethodGet, nil, nil, nil, nil)

// 			req := httptest.NewRequest("GET", "/test", nil)
// 			req.RemoteAddr = tc.clientID
// 			ginCtx.Request = req

// 			s := &sse{
// 				clients: make(map[string][]*SSEConn),
// 			}

// 			_, err := s.HandleConnection(ginCtx)
// 			if (err != nil) != tc.wantErr {
// 				t.Errorf("HandleConnection() error = %v, wantErr %v", err, tc.wantErr)
// 			}

// 			if !tc.wantErr {
// 				s.BroadcastMessage("test message")
// 				closeChannel <- true
// 				require.Equal(t, tc.message, w.Body.String())
// 			}
// 		})
// 	}
// }

func Test_sse_BroadcastMessage(t *testing.T) {
	// Create a new SSE server
	s := &sse{
		clients: make(map[string][]*SSEConn),
		mutex:   sync.Mutex{},
	}

	// Register a client
	clientID := "client1"
	messageChannel := make(chan string)
	s.clients[clientID] = []*SSEConn{
		{
			Channel: messageChannel,
			ID:      clientID,
		},
	}

	message := "test message"

	go func() {
		// Broadcast a message
		err := s.BroadcastMessage(message)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}()

	// Check that the message was received
	receivedMessage := <-messageChannel
	if receivedMessage != message {
		t.Errorf("Expected message '%s', got '%s'", message, receivedMessage)

	}
}

func Test_sse_SendMessage(t *testing.T) {
	s := &sse{
		clients: make(map[string][]*SSEConn),
		mutex:   sync.Mutex{},
	}

	// Create a dummy client
	s.clients["dummyUser"] = []*SSEConn{
		{
			Channel: make(chan string, 1),
			ID:      "dummyUser-client1",
		},
	}

	tests := map[string]struct {
		userID  string
		message string
		wantErr bool
	}{
		"success": {
			userID:  "dummyUser",
			message: "Hello, World!",
			wantErr: false,
		},
		"failure": {
			userID:  "nonExistentUser",
			message: "Hello, World!",
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := s.SendMessage(tc.userID, tc.message)

			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got %v", err)
				}

				select {
				case msg := <-s.clients[tc.userID][0].Channel:
					if msg != tc.message {
						t.Errorf("expected message %q but got %q", tc.message, msg)
					}
				default:
					t.Errorf("expected message %q but got none", tc.message)
				}
			}
		})
	}

}

func Test_sse_SendData(t *testing.T) {
	type testStruct struct {
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
	s := &sse{
		clients: make(map[string][]*SSEConn),
		mutex:   sync.Mutex{},
	}
	// Create a dummy client
	s.clients["user1"] = []*SSEConn{
		&SSEConn{
			Channel: make(chan string, 1),
			ID:      "user1-client1",
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
			err := s.SendData(tc.userID, tc.data)
			if (err != nil) != tc.wantErr {
				t.Fatalf("%v case: SendData() error = %v, wantErr %v", name, err, tc.wantErr)
			}
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got %v", err)
				}

				select {
				case msg := <-s.clients[tc.userID][0].Channel:
					if msg != tc.message {
						t.Errorf("expected message %q but got %q", tc.message, msg)
					}
				default:
					t.Errorf("expected message %q but got none", tc.message)
				}
			}
		})
	}
}

func Test_sse_BroadcastData(t *testing.T) {
	type testStruct struct {
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}
	s := &sse{
		clients: make(map[string][]*SSEConn),
		mutex:   sync.Mutex{},
	}
	// Create a dummy client
	dummyClient := make(chan string, 1)
	s.clients["user1"] = []*SSEConn{&SSEConn{
		Channel: dummyClient,
		ID:      "user1-client1",
	}}

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
			err := s.BroadcastData(tc.data)
			if (err != nil) != tc.wantErr {
				t.Fatalf("%v case: BroadcastData() error = %v, wantErr %v", name, err, tc.wantErr)
			}
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got %v", err)
				}

				select {
				case msg := <-dummyClient:
					if msg != tc.message {
						t.Errorf("expected message %q but got %q", tc.message, msg)
					}
				default:
					t.Errorf("expected message %q but got none", tc.message)
				}
			}
		})
	}
}
