package user

import (
	"context"
	"reflect"
	"testing"

	mocks "github.com/dwarvesf/go-api/mocks/pkg/repository/user"
	"github.com/dwarvesf/go-api/pkg/config"
	"github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_impl_UpdateUser(t *testing.T) {
	type mocked struct {
		uID                 int
		expGetUserCalled    bool
		getUser             *model.User
		getUserErr          error
		expUpdateUserCalled bool
		updateUser          *model.User
		updateUserErr       error
	}
	type args struct {
		req  model.UpdateUserRequest
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
				expUpdateUserCalled: true,
				updateUser: &model.User{
					ID:             1,
					Email:          "admin@d.foundation",
					FullName:       "admin1",
					Status:         "active",
					Avatar:         "https://d.foundation/avatar1.png",
					Role:           "admin",
					HashedPassword: "hash",
					Salt:           "abcdef",
				},
			},
			args: args{
				req: model.UpdateUserRequest{
					FullName: "admin1",
					Avatar:   "https://d.foundation/avatar1.png",
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
				uID:                 2,
				expGetUserCalled:    true,
				getUserErr:          model.ErrNotFound,
				expUpdateUserCalled: false,
			},
			args: args{
				req: model.UpdateUserRequest{
					FullName: "admin1",
					Avatar:   "https://d.foundation/avatar1.png",
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

			if tt.mocked.expUpdateUserCalled {
				userRepoMock.
					EXPECT().
					Update(mock.Anything, mock.Anything, mock.Anything).
					Return(tt.mocked.updateUser, tt.mocked.updateUserErr)
			}

			c := &impl{
				repo: &repository.Repo{
					User: userRepoMock,
				},
				cfg: config.LoadTestConfig(),
			}

			_, err := db.Init(c.cfg)
			require.NoError(t, err)

			ctx := context.Background()
			// assign userID to context
			ctx = context.WithValue(ctx, middleware.UserIDCtxKey, tt.mocked.uID)
			got, err := c.UpdateUser(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("impl.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("impl.UpdateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
