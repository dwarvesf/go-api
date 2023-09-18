package user

import (
	"context"
	"testing"

	mocks "github.com/dwarvesf/go-api/mocks/pkg/repository/user"
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/monitor"
	"github.com/dwarvesf/go-api/pkg/repository"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_impl_UpdatePassword(t *testing.T) {
	type mocked struct {
		uID                     int
		expGetUserCalled        bool
		getUser                 *model.User
		getUserErr              error
		expUpdatePasswordCalled bool
		updatePasswordErr       error
	}
	type args struct {
		req  model.UpdatePasswordRequest
		role string
	}
	tests := map[string]struct {
		mocked  mocked
		args    args
		want    *model.User
		wantErr bool
	}{
		"success": {
			mocked: mocked{
				uID:              1,
				expGetUserCalled: true,
				getUser: &model.User{
					ID:             1,
					Email:          "admin@d.foundation",
					FullName:       "admin",
					Status:         "active",
					Avatar:         "https://d.foundation/avatar.png",
					Role:           "admin",
					HashedPassword: "hash",
					Salt:           "abcdef",
				},
				expUpdatePasswordCalled: true,
			},
			args: args{
				req: model.UpdatePasswordRequest{
					OldPassword: "hash",
					NewPassword: "hash1",
				},
			},
			want: &model.User{
				ID:             1,
				Email:          "admin@d.foundation",
				FullName:       "admin1",
				Status:         "active",
				Avatar:         "https://d.foundation/avatar1.png",
				Role:           "admin",
				HashedPassword: "hash",
				Salt:           "abcdef",
			},
			wantErr: false,
		},
		"not found": {
			mocked: mocked{
				uID:                     2,
				expGetUserCalled:        true,
				getUserErr:              model.ErrNotFound,
				expUpdatePasswordCalled: false,
			},
			args: args{
				req: model.UpdatePasswordRequest{
					OldPassword: "hash",
					NewPassword: "hash1",
				},
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				userRepoMock = mocks.NewRepo(t)
			)

			if tt.mocked.expGetUserCalled {
				userRepoMock.
					EXPECT().
					GetByID(mock.Anything, mock.Anything).
					Return(tt.mocked.getUser, tt.mocked.getUserErr)
			}

			if tt.mocked.expUpdatePasswordCalled {
				userRepoMock.
					EXPECT().
					UpdatePassword(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mocked.updatePasswordErr)
			}

			c := &impl{
				repo: &repository.Repo{
					User: userRepoMock,
				},
				cfg:     config.LoadTestConfig(),
				monitor: monitor.TestMonitor(),
			}

			_, err := db.Init(c.cfg)
			require.NoError(t, err)

			ctx := context.Background()
			// assign userID to context
			ctx = context.WithValue(ctx, middleware.UserIDCtxKey, tt.mocked.uID)
			err = c.UpdatePassword(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("impl.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
