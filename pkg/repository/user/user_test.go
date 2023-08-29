package user

import (
	"reflect"
	"testing"

	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func Test_repo_GetByID(t *testing.T) {
	db.WithTestingDB(t, func(ctx db.Context) {
		u := &orm.User{
			Email:          "admin@d.foundation",
			Name:           "admin",
			Status:         "active",
			Avatar:         "https://d.foundation/avatar.png",
			Role:           "admin",
			HashedPassword: "123456",
			Salt:           "abcdef",
		}
		err := u.Insert(ctx, ctx.DB, boil.Infer())
		require.NoError(t, err)

		type args struct {
			uID int
		}
		tests := map[string]struct {
			args    args
			want    *model.User
			wantErr bool
		}{
			"success": {
				args: args{
					uID: u.ID,
				},
				want: &model.User{
					ID:             u.ID,
					Email:          u.Email,
					FullName:       u.Name,
					Status:         u.Status,
					Avatar:         u.Avatar,
					HashedPassword: u.HashedPassword,
					Role:           u.Role,
				},
				wantErr: false,
			},
			"not found": {
				args: args{
					uID: u.ID + 1,
				},
				want:    nil,
				wantErr: true,
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				r := &repo{}
				got, err := r.GetByID(ctx, tt.args.uID)
				if (err != nil) != tt.wantErr {
					t.Errorf("repo.GetByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("repo.GetByID() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
