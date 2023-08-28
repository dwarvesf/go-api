package portal

import (
	"net/http/httptest"
	"testing"

	mocks "github.com/dwarvesf/go-api/mocks/pkg/controller/user"
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/handler/testutil"
	"github.com/dwarvesf/go-api/pkg/logger"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Me(t *testing.T) {
	type mocked struct {
		expUpdateJWT bool
		userID       int
		role         string
		expGetUser   bool
		user         *orm.User
		userErr      error
	}

	type expected struct {
		Status int
		Body   string
	}
	tests := []struct {
		name     string
		mocked   mocked
		expected expected
	}{
		{
			name: "success",
			mocked: mocked{
				expUpdateJWT: true,
				userID:       1,
				expGetUser:   true,
				user: &orm.User{
					ID:    "1",
					Email: "admin@email.com",
				},
			},
			expected: expected{
				Status: 200,
				Body:   "admin",
			},
		},
	}
	for _, tt := range tests {
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
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				log:      logger.NewLogger(),
				cfg:      cfg,
				userCtrl: ctrlMock,
			}
			h.Me(ginCtx)

			assert.Equal(t, tt.expected.Status, w.Code)
			resBody := w.Body.String()
			assert.Contains(t, resBody, tt.expected.Body)
		})
	}
}
