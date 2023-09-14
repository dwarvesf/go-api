package auth

import (
	"database/sql"
	"errors"
	"net/http/httptest"
	"reflect"
	"testing"

	mocks "github.com/dwarvesf/go-api/mocks/pkg/repository/user"
	jwtmocks "github.com/dwarvesf/go-api/mocks/pkg/service/jwthelper"
	passworkmocks "github.com/dwarvesf/go-api/mocks/pkg/service/passwordhelper"
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/handler/testutil"
	"github.com/dwarvesf/go-api/pkg/logger/monitor"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/dwarvesf/go-api/pkg/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_impl_Login(t *testing.T) {
	validPass, err := util.GenerateHashedKey("123456")
	require.NoError(t, err)

	type mocked struct {
		expGetUserCalled bool
		getUser          *model.User
		getUserErr       error
		expJWTCalled     bool
		jwtToken         string
		genjwtErr        error
		compareCalled    bool
		compare          bool
	}
	type args struct {
		req  model.LoginRequest
		role string
	}
	tests := map[string]struct {
		mocked  mocked
		args    args
		want    *model.LoginResponse
		wantErr bool
	}{
		"success": {
			mocked: mocked{
				expGetUserCalled: true,
				getUser: &model.User{
					ID:             1,
					Email:          "admin@d.foundation",
					FullName:       "admin",
					Status:         "active",
					Avatar:         "https://d.foundation/avatar.png",
					Role:           "admin",
					HashedPassword: validPass,
					Salt:           "abcdef",
				},
				expJWTCalled:  true,
				jwtToken:      "token",
				compareCalled: true,
				compare:       true,
			},
			args: args{
				req: model.LoginRequest{
					Email:    "admin@d.foundation",
					Password: "123456",
				},
				role: "admin",
			},
			want: &model.LoginResponse{
				ID:          1,
				Email:       "admin@d.foundation",
				AccessToken: "token",
			},
			wantErr: false,
		},
		"invalid password": {
			mocked: mocked{
				expGetUserCalled: true,
				getUser: &model.User{
					ID:             1,
					Email:          "admin@d.foundation",
					FullName:       "admin",
					Status:         "active",
					Avatar:         "https://d.foundation/avatar.png",
					Role:           "admin",
					HashedPassword: validPass,
					Salt:           "abcdef",
				},
				expJWTCalled:  false,
				jwtToken:      "token",
				compareCalled: true,
				compare:       false,
			},
			args: args{
				req: model.LoginRequest{
					Email:    "admin@d.foundation",
					Password: "123458",
				},
				role: "admin",
			},
			want:    nil,
			wantErr: true,
		},
		"not found user": {
			mocked: mocked{
				expGetUserCalled: true,
				getUser: &model.User{
					ID:             1,
					Email:          "admin@d.foundation",
					FullName:       "admin",
					Status:         "active",
					Avatar:         "https://d.foundation/avatar.png",
					Role:           "admin",
					HashedPassword: validPass,
					Salt:           "abcdef",
				},
				getUserErr:    sql.ErrNoRows,
				expJWTCalled:  false,
				jwtToken:      "token",
				compareCalled: false,
				compare:       false,
			},
			args: args{
				req: model.LoginRequest{
					Email:    "admin@d.foundation",
					Password: "123458",
				},
				role: "admin",
			},
			want:    nil,
			wantErr: true,
		},
		"generate token failed": {
			mocked: mocked{
				expGetUserCalled: true,
				getUser: &model.User{
					ID:             1,
					Email:          "admin@d.foundation",
					FullName:       "admin",
					Status:         "active",
					Avatar:         "https://d.foundation/avatar.png",
					Role:           "admin",
					HashedPassword: validPass,
					Salt:           "abcdef",
				},
				expJWTCalled:  true,
				genjwtErr:     errors.New("failed to generate token"),
				compareCalled: true,
				compare:       true,
			},
			args: args{
				req: model.LoginRequest{
					Email:    "admin@d.foundation",
					Password: "123458",
				},
				role: "admin",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				userRepoMock = mocks.NewRepo(t)
				jwtMock      = jwtmocks.NewHelper(t)
				passwordMock = passworkmocks.NewHelper(t)
			)

			if tt.mocked.expGetUserCalled {
				userRepoMock.
					EXPECT().
					GetByEmail(mock.Anything, mock.Anything).
					Return(tt.mocked.getUser, tt.mocked.getUserErr)
			}

			if tt.mocked.compareCalled {
				passwordMock.
					EXPECT().
					Compare(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mocked.compare)
			}

			if tt.mocked.expJWTCalled {
				jwtMock.
					EXPECT().
					GenerateJWTToken(mock.Anything).
					Return(tt.mocked.jwtToken, tt.mocked.genjwtErr)
			}
			c := &impl{
				repo: &repository.Repo{
					User: userRepoMock,
				},
				jwtHelper:      jwtMock,
				passwordHelper: passwordMock,
				cfg:            config.LoadTestConfig(),
				monitor:        monitor.TestMonitor(),
			}

			_, err = db.Init(c.cfg)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			ginCtx := testutil.NewRequest(w, testutil.MethodPost, nil, nil, nil, nil)
			got, err := c.Login(ginCtx.Request.Context(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("impl.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("impl.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
