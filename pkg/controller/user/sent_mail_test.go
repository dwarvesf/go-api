package user

import (
	"context"
	"errors"
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

func Test_impl_SentMail(t *testing.T) {
	type mocked struct {
		uID              int
		expGetListCalled bool
		users            *model.ListResult[model.User]
		GetListErr       error
		countCalled      bool
		countErr         error
		count            int64
	}
	tests := map[string]struct {
		mocked  mocked
		want    *model.User
		wantErr bool
	}{
		"success": {
			mocked: mocked{
				uID:              1,
				expGetListCalled: true,
				users: &model.ListResult[model.User]{
					Pagination: model.Pagination{
						Page:         1,
						PageSize:     10,
						TotalRecords: 1,
						TotalPages:   1,
						Offset:       0,
						Sort:         "",
						HasNext:      false,
					},
					Data: []model.User{
						{
							ID:             1,
							Email:          "admin@d.foundation",
							FullName:       "admin",
							Status:         "active",
							Avatar:         "https://d.foundation/avatar.png",
							Role:           "admin",
							HashedPassword: "hash",
							Salt:           "abcdef",
						},
					},
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
		"failed to get list": {
			mocked: mocked{
				uID:              1,
				expGetListCalled: true,
				GetListErr:       errors.New("failed to get list"),
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var (
				userRepoMock = mocks.NewRepo(t)
			)

			if tt.mocked.countCalled {
				userRepoMock.
					EXPECT().
					Count(mock.Anything).
					Return(tt.mocked.count, tt.mocked.countErr)
			}

			if tt.mocked.expGetListCalled {
				userRepoMock.
					EXPECT().
					GetList(mock.Anything, mock.Anything).
					Return(tt.mocked.users, tt.mocked.GetListErr)
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
			err = c.SentMail(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("impl.SentMail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
