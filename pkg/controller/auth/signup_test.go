package auth

import (
	"database/sql"
	"net/http/httptest"
	"testing"

	mocks "github.com/dwarvesf/go-api/mocks/pkg/repository/user"
	jwtmocks "github.com/dwarvesf/go-api/mocks/pkg/service/jwthelper"
	passworkmocks "github.com/dwarvesf/go-api/mocks/pkg/service/passwordhelper"
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/handler/testutil"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_impl_Signup(t *testing.T) {
	type mocked struct {
		expGetUserCalled    bool
		getUser             *model.User
		getUserErr          error
		expCreateUserCalled bool
		createUser          *model.User
		createUserErr       error
		genSaltCalled       bool
		salt                string
		hashCalled          bool
		hash                string
		hashErr             error
	}
	type args struct {
		req  model.SignupRequest
		role string
	}
	tests := map[string]struct {
		mocked  mocked
		args    args
		wantErr bool
	}{
		"success": {
			mocked: mocked{
				genSaltCalled:       true,
				salt:                "salt",
				hashCalled:          true,
				hash:                "hash",
				expGetUserCalled:    true,
				getUserErr:          sql.ErrNoRows,
				expCreateUserCalled: true,
				createUser: &model.User{
					Email:          "admin@d.foundation",
					FullName:       "admin",
					Status:         "active",
					Avatar:         "https://d.foundation/avatar.png",
					Role:           "admin",
					HashedPassword: "hash",
					Salt:           "abcdef",
				},
			},
			args: args{
				req: model.SignupRequest{
					Email:    "admin@d.foundation",
					Password: "123456",
					Name:     "admin",
					Avatar:   "https://d.foundation/avatar.png",
				},
				role: "admin",
			},
			wantErr: false,
		},
		"duplicate email": {
			mocked: mocked{
				genSaltCalled:       true,
				salt:                "salt",
				hashCalled:          true,
				hash:                "hash",
				expGetUserCalled:    true,
				expCreateUserCalled: false,
				getUser: &model.User{
					Email:          "admin@d.foundation",
					FullName:       "admin",
					Status:         "active",
					Avatar:         "https://d.foundation/avatar.png",
					Role:           "admin",
					HashedPassword: "hash",
					Salt:           "abcdef",
				},
			},
			args: args{
				req: model.SignupRequest{
					Email:    "admin@d.foundation",
					Password: "123456",
					Name:     "admin",
					Avatar:   "https://d.foundation/avatar.png",
				},
				role: "admin",
			},
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

			if tt.mocked.genSaltCalled {
				passwordMock.
					EXPECT().
					GenerateSalt().
					Return(tt.mocked.salt)
			}

			if tt.mocked.hashCalled {
				passwordMock.
					EXPECT().
					Hash(mock.Anything, mock.Anything).
					Return(tt.mocked.hash, tt.mocked.hashErr)
			}

			if tt.mocked.expGetUserCalled {
				userRepoMock.
					EXPECT().
					GetByEmail(mock.Anything, mock.Anything).
					Return(tt.mocked.getUser, tt.mocked.getUserErr)
			}

			if tt.mocked.expCreateUserCalled {
				userRepoMock.
					EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(tt.mocked.createUser, tt.mocked.createUserErr)
			}

			c := &impl{
				repo: &repository.Repo{
					User: userRepoMock,
				},
				jwtHelper:      jwtMock,
				passwordHelper: passwordMock,
				cfg:            config.LoadTestConfig(),
			}

			_, err := db.Init(c.cfg)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			ginCtx := testutil.NewRequest(w, testutil.MethodPost, nil, nil, nil, nil)
			err = c.Signup(ginCtx.Request.Context(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("impl.Signup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
