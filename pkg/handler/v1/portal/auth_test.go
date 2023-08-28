package portal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/dwarvesf/go-api/mocks/pkg/controller/auth"
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/handler/testutil"
	"github.com/dwarvesf/go-api/pkg/handler/v1/view"
	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Login(t *testing.T) {
	type mocked struct {
		expLoginCalled bool
		loginResponse  *model.LoginResponse
		loginErr       error
	}
	type args struct {
		input view.LoginRequest
	}
	type expected struct {
		Status int
		Body   string
	}

	tests := []struct {
		name     string
		mocked   mocked
		args     args
		expected expected
	}{
		{
			name: "success",
			mocked: mocked{
				expLoginCalled: true,
				loginResponse: &model.LoginResponse{
					ID:          "1",
					Email:       "admin@email.com",
					AccessToken: "token",
				},
				loginErr: nil,
			},
			args: args{
				input: view.LoginRequest{
					Email:    "admin@gmail.com",
					Password: "abcd1234",
				},
			},
			expected: expected{
				Status: http.StatusOK,
				Body:   "token",
			},
		},
		{
			name: "error",
			mocked: mocked{
				expLoginCalled: true,
				loginResponse:  nil,
				loginErr:       model.ErrInvalidCredentials,
			},
			args: args{
				input: view.LoginRequest{
					Email:    "admin@gmail.com",
					Password: "invalid",
				},
			},
			expected: expected{
				Status: http.StatusBadRequest,
				Body:   "Wrong username or password",
			},
		},
		{
			name: "bad request",
			mocked: mocked{
				expLoginCalled: false,
			},
			args: args{
				input: view.LoginRequest{
					Email:    "admin@gmail.com",
					Password: "",
				},
			},
			expected: expected{
				Status: http.StatusBadRequest,
				Body:   "required",
			},
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		ginCtx := testutil.NewRequest(w, testutil.MethodPost, nil, nil, nil, tt.args.input)

		var (
			ctrlMock = mocks.NewController(t)
		)

		if tt.mocked.expLoginCalled {
			ctrlMock.EXPECT().Login(mock.Anything).Return(tt.mocked.loginResponse, tt.mocked.loginErr)
		}
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				log:      logger.NewLogger(),
				cfg:      config.LoadTestConfig(),
				authCtrl: ctrlMock,
			}
			h.Login(ginCtx)

			assert.Equal(t, tt.expected.Status, w.Code)
			resBody := w.Body.String()
			assert.Contains(t, resBody, tt.expected.Body)
		})
	}
}

func TestHandler_Signup(t *testing.T) {
	type mocked struct {
		expSignupCalled bool
		signupErr       error
	}
	type args struct {
		input view.SignupRequest
	}
	type expected struct {
		Status int
		Body   view.MessageResponse
	}

	tests := []struct {
		name     string
		mocked   mocked
		args     args
		expected expected
	}{
		{
			name: "success",
			mocked: mocked{
				expSignupCalled: true,
				signupErr:       nil,
			},
			args: args{
				input: view.SignupRequest{
					Email:    "admin@gmail.com",
					Password: "abcd1234",
					FullName: "Admin",
					Status:   "active",
					Avatar:   "https://www.google.com",
				},
			},
			expected: expected{
				Status: http.StatusCreated,
				Body: view.MessageResponse{
					Data: view.Message{
						Message: "success",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		ginCtx := testutil.NewRequest(w, testutil.MethodPost, nil, nil, nil, tt.args.input)

		var (
			ctrlMock = mocks.NewController(t)
		)

		if tt.mocked.expSignupCalled {
			ctrlMock.EXPECT().Signup(mock.Anything).Return(tt.mocked.signupErr)
		}
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				log:      logger.NewLogger(),
				cfg:      config.LoadTestConfig(),
				authCtrl: ctrlMock,
			}
			h.Signup(ginCtx)
			assert.Equal(t, tt.expected.Status, w.Code)
			resBody := w.Body.String()
			body, err := json.Marshal(tt.expected.Body)
			assert.Nil(t, err)
			assert.Equal(t, resBody, string(body))
		})
	}
}
