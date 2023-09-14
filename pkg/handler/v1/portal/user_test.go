package portal

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	mocks "github.com/dwarvesf/go-api/mocks/pkg/controller/user"
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/handler/testutil"
	"github.com/dwarvesf/go-api/pkg/handler/v1/view"
	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/dwarvesf/go-api/pkg/logger/monitor"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Me(t *testing.T) {
	type mocked struct {
		expUpdateJWT bool
		userID       int
		role         string
		expGetUser   bool
		user         *model.User
		userErr      error
	}

	type expected struct {
		Status int
		Body   string
	}
	tests := map[string]struct {
		mocked   mocked
		expected expected
	}{
		"success": {
			mocked: mocked{
				expUpdateJWT: true,
				userID:       1,
				expGetUser:   true,
				user: &model.User{
					ID:    1,
					Email: "admin@email.com",
				},
			},
			expected: expected{
				Status: 200,
				Body:   "admin",
			},
		},
	}
	for name, tt := range tests {
		w := httptest.NewRecorder()
		cfg := config.LoadTestConfig()
		ginCtx := testutil.NewRequest(w, testutil.MethodGet, nil, nil, nil, nil)

		if tt.mocked.expUpdateJWT {
			testutil.UpdateJWT(ginCtx, tt.mocked.userID, tt.mocked.role)
		}

		var (
			ctrlMock = mocks.NewController(t)
		)

		if tt.mocked.expGetUser {
			ctrlMock.EXPECT().Me(mock.Anything).Return(tt.mocked.user, tt.mocked.userErr)
		}
		t.Run(name, func(t *testing.T) {
			h := Handler{
				log:      logger.NewLogger(),
				cfg:      cfg,
				userCtrl: ctrlMock,
				monitor:  monitor.TestMonitor(),
			}
			h.Me(ginCtx)

			assert.Equal(t, tt.expected.Status, w.Code)
			resBody := w.Body.String()
			assert.Contains(t, resBody, tt.expected.Body)
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	type mocked struct {
		expUpdateJWT  bool
		userID        int
		role          string
		expUpdateUser bool
		user          *model.User
		userErr       error
	}

	type args struct {
		input view.UpdateUserRequest
	}

	type expected struct {
		Status int
		Body   view.UserResponse
	}
	tests := map[string]struct {
		args     args
		mocked   mocked
		expected expected
	}{
		"success": {
			mocked: mocked{
				expUpdateJWT:  true,
				userID:        1,
				expUpdateUser: true,
				user: &model.User{
					ID:    1,
					Email: "admin@email.com",
				},
			},
			args: args{
				input: view.UpdateUserRequest{
					FullName: "Admin",
					Avatar:   "https://www.google.com",
				},
			},
			expected: expected{
				Status: 200,
				Body: view.UserResponse{
					Data: view.User{
						ID:    1,
						Email: "admin@email.com",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		w := httptest.NewRecorder()
		cfg := config.LoadTestConfig()
		ginCtx := testutil.NewRequest(w, testutil.MethodPut, nil, nil, nil, tt.args.input)

		if tt.mocked.expUpdateJWT {
			testutil.UpdateJWT(ginCtx, tt.mocked.userID, tt.mocked.role)
		}

		var (
			ctrlMock = mocks.NewController(t)
		)

		if tt.mocked.expUpdateUser {
			ctrlMock.EXPECT().UpdateUser(mock.Anything, mock.Anything).Return(tt.mocked.user, tt.mocked.userErr)
		}
		t.Run(name, func(t *testing.T) {
			h := Handler{
				log:      logger.NewLogger(),
				cfg:      cfg,
				userCtrl: ctrlMock,
				monitor:  monitor.TestMonitor(),
			}
			h.UpdateUser(ginCtx)

			assert.Equal(t, tt.expected.Status, w.Code)
			resBody := w.Body.String()
			body, err := json.Marshal(tt.expected.Body)
			assert.Nil(t, err)
			assert.Equal(t, resBody, string(body))
		})
	}
}

func TestHandler_UpdatePassword(t *testing.T) {
	type mocked struct {
		expUpdateJWT      bool
		userID            int
		role              string
		expUpdatePassword bool
		userErr           error
	}

	type args struct {
		input view.UpdatePasswordRequest
	}

	type expected struct {
		Status int
		Body   view.MessageResponse
	}
	tests := map[string]struct {
		args     args
		mocked   mocked
		expected expected
	}{
		"success": {
			mocked: mocked{
				expUpdateJWT:      true,
				userID:            1,
				expUpdatePassword: true,
			},
			args: args{
				input: view.UpdatePasswordRequest{
					NewPassword: "123456",
					OldPassword: "123456",
				},
			},
			expected: expected{
				Status: 200,
				Body: view.MessageResponse{
					Data: view.Message{
						Message: "success",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		w := httptest.NewRecorder()
		cfg := config.LoadTestConfig()
		ginCtx := testutil.NewRequest(w, testutil.MethodPut, nil, nil, nil, tt.args.input)

		if tt.mocked.expUpdateJWT {
			testutil.UpdateJWT(ginCtx, tt.mocked.userID, tt.mocked.role)
		}

		var (
			ctrlMock = mocks.NewController(t)
		)

		if tt.mocked.expUpdatePassword {
			ctrlMock.EXPECT().UpdatePassword(mock.Anything, mock.Anything).Return(tt.mocked.userErr)
		}
		t.Run(name, func(t *testing.T) {
			h := Handler{
				log:      logger.NewLogger(),
				cfg:      cfg,
				userCtrl: ctrlMock,
				monitor:  monitor.TestMonitor(),
			}
			h.UpdatePassword(ginCtx)

			assert.Equal(t, tt.expected.Status, w.Code)
			resBody := w.Body.String()
			body, err := json.Marshal(tt.expected.Body)
			assert.Nil(t, err)
			assert.Equal(t, resBody, string(body))
		})
	}
}
