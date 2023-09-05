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
					Salt:           u.Salt,
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
				require.Equal(t, tt.want, got)
			})
		}
	})
}

func Test_repo_Count(t *testing.T) {
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

		tests := map[string]struct {
			want    int64
			wantErr bool
		}{
			"success": {
				want:    1,
				wantErr: false,
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				r := &repo{}
				got, err := r.Count(ctx)
				if (err != nil) != tt.wantErr {
					t.Errorf("repo.Count() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("repo.Count() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func Test_repo_GetByEmail(t *testing.T) {
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
			email string
		}
		tests := map[string]struct {
			args    args
			want    *model.User
			wantErr bool
		}{
			"success": {
				args: args{
					email: u.Email,
				},
				want: &model.User{
					ID:             u.ID,
					Email:          u.Email,
					FullName:       u.Name,
					Status:         u.Status,
					Avatar:         u.Avatar,
					HashedPassword: u.HashedPassword,
					Role:           u.Role,
					Salt:           u.Salt,
				},
				wantErr: false,
			},
			"not found": {
				args: args{
					email: "a@gmail.com",
				},
				want:    nil,
				wantErr: true,
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				r := &repo{}
				got, err := r.GetByEmail(ctx, tt.args.email)
				if (err != nil) != tt.wantErr {
					t.Errorf("repo.GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("repo.GetByEmail() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func Test_repo_Create(t *testing.T) {
	db.WithTestingDB(t, func(ctx db.Context) {
		u := &orm.User{
			Email:          "admin1@d.foundation",
			Name:           "admin1",
			Status:         "active",
			Avatar:         "https://d.foundation/avatar.png",
			Role:           "admin",
			HashedPassword: "123456",
			Salt:           "abcdef",
		}
		err := u.Insert(ctx, ctx.DB, boil.Infer())
		require.NoError(t, err)

		type args struct {
			req model.SignupRequest
		}
		tests := map[string]struct {
			args    args
			want    *model.User
			wantErr bool
		}{
			"success": {
				args: args{
					req: model.SignupRequest{
						Email:          "admin@d.foundation",
						Name:           "admin",
						Status:         "active",
						Avatar:         "https://d.foundation/avatar.png",
						Role:           "admin",
						HashedPassword: "123456",
						Salt:           "abcdef",
					},
				},
				want: &model.User{
					Email:          "admin@d.foundation",
					FullName:       "admin",
					Status:         "active",
					Avatar:         "https://d.foundation/avatar.png",
					Role:           "admin",
					HashedPassword: "123456",
					Salt:           "abcdef",
				},
				wantErr: false,
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				r := &repo{}
				got, err := r.Create(ctx, tt.args.req)
				if (err != nil) != tt.wantErr {
					t.Errorf("repo.Create() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !tt.wantErr {
					got.ID = 0
					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf("repo.Create() = %v, want %v", got, tt.want)
					}
				}
			})
		}
	})
}

func Test_repo_Update(t *testing.T) {
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
			uID  int
			user model.UpdateUserRequest
		}
		tests := map[string]struct {
			args    args
			want    *model.User
			wantErr bool
		}{
			"success": {
				args: args{
					uID: u.ID,
					user: model.UpdateUserRequest{
						FullName: "admin1",
						Avatar:   "https://d.foundation/avatar2.png",
					},
				},
				want: &model.User{
					ID:             u.ID,
					Email:          u.Email,
					FullName:       "admin1",
					Avatar:         "https://d.foundation/avatar2.png",
					Status:         u.Status,
					HashedPassword: u.HashedPassword,
					Role:           u.Role,
					Salt:           u.Salt,
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
				got, err := r.Update(ctx, tt.args.uID, tt.args.user)
				if (err != nil) != tt.wantErr {
					t.Errorf("repo.Update() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("repo.Update() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func Test_repo_UpdatePassword(t *testing.T) {
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
			uID         int
			newPassword string
		}
		tests := map[string]struct {
			args    args
			want    *model.User
			wantErr bool
		}{
			"success": {
				args: args{
					uID:         u.ID,
					newPassword: "1234567",
				},
				want: &model.User{
					ID:             u.ID,
					Email:          u.Email,
					FullName:       u.Name,
					Status:         u.Status,
					Avatar:         u.Avatar,
					HashedPassword: "1234567",
					Role:           u.Role,
					Salt:           u.Salt,
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
				err := r.UpdatePassword(ctx, tt.args.uID, tt.args.newPassword)
				if (err != nil) != tt.wantErr {
					t.Errorf("repo.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			})
		}
	})
}

func Test_repo_GetList(t *testing.T) {
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

		u = &orm.User{
			Email:          "admin1@d.foundation",
			Name:           "admin1",
			Status:         "active",
			Avatar:         "https://d.foundation/avatar1.png",
			Role:           "admin",
			HashedPassword: "123456",
			Salt:           "abcdef",
		}
		err = u.Insert(ctx, ctx.DB, boil.Infer())
		require.NoError(t, err)

		type args struct {
			page     int
			pageSize int
			sort     string
			query    string
		}
		tests := map[string]struct {
			args    args
			want    *model.UserList
			wantErr bool
		}{
			"success": {
				args: args{
					page:     1,
					pageSize: 10,
					sort:     "+id",
					query:    "",
				},
				want: &model.UserList{
					Pagination: model.Pagination{
						Page:         1,
						PageSize:     10,
						TotalRecords: 2,
						TotalPages:   1,
						Offset:       0,
						Sort:         "+id",
						HasNext:      false,
					},
					Data: []*model.User{
						{
							Email:          "admin@d.foundation",
							FullName:       "admin",
							Status:         "active",
							Avatar:         "https://d.foundation/avatar.png",
							Role:           "admin",
							HashedPassword: "123456",
							Salt:           "abcdef",
						},
						{
							Email:          "admin1@d.foundation",
							FullName:       "admin1",
							Status:         "active",
							Avatar:         "https://d.foundation/avatar1.png",
							Role:           "admin",
							HashedPassword: "123456",
							Salt:           "abcdef",
						},
					},
				},
				wantErr: false,
			},
			"with sort": {
				args: args{
					page:     1,
					pageSize: 10,
					sort:     "-id",
					query:    "",
				},
				want: &model.UserList{
					Pagination: model.Pagination{
						Page:         1,
						PageSize:     10,
						TotalRecords: 2,
						TotalPages:   1,
						Offset:       0,
						Sort:         "-id",
						HasNext:      false,
					},
					Data: []*model.User{
						{
							Email:          "admin1@d.foundation",
							FullName:       "admin1",
							Status:         "active",
							Avatar:         "https://d.foundation/avatar1.png",
							Role:           "admin",
							HashedPassword: "123456",
							Salt:           "abcdef",
						},
						{
							Email:          "admin@d.foundation",
							FullName:       "admin",
							Status:         "active",
							Avatar:         "https://d.foundation/avatar.png",
							Role:           "admin",
							HashedPassword: "123456",
							Salt:           "abcdef",
						},
					},
				},
				wantErr: false,
			},
			"with query": {
				args: args{
					page:     1,
					pageSize: 10,
					sort:     "-id",
					query:    "admin",
				},
				want: &model.UserList{
					Pagination: model.Pagination{
						Page:         1,
						PageSize:     10,
						TotalRecords: 2,
						TotalPages:   1,
						Offset:       0,
						Sort:         "-id",
						HasNext:      false,
					},
					Data: []*model.User{
						{
							Email:          "admin1@d.foundation",
							FullName:       "admin1",
							Status:         "active",
							Avatar:         "https://d.foundation/avatar1.png",
							Role:           "admin",
							HashedPassword: "123456",
							Salt:           "abcdef",
						},
						{
							Email:          "admin@d.foundation",
							FullName:       "admin",
							Status:         "active",
							Avatar:         "https://d.foundation/avatar.png",
							Role:           "admin",
							HashedPassword: "123456",
							Salt:           "abcdef",
						},
					},
				},
				wantErr: false,
			},
			"with page": {
				args: args{
					page:     2,
					pageSize: 1,
					sort:     "-id",
					query:    "",
				},
				want: &model.UserList{
					Pagination: model.Pagination{
						Page:         2,
						PageSize:     1,
						TotalRecords: 2,
						TotalPages:   2,
						Offset:       1,
						Sort:         "-id",
						HasNext:      false,
					},
					Data: []*model.User{
						{
							Email:          "admin@d.foundation",
							FullName:       "admin",
							Status:         "active",
							Avatar:         "https://d.foundation/avatar.png",
							Role:           "admin",
							HashedPassword: "123456",
							Salt:           "abcdef",
						},
					},
				},
				wantErr: false,
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				r := &repo{}
				got, err := r.GetList(ctx, tt.args.page, tt.args.pageSize, tt.args.sort, tt.args.query)
				if (err != nil) != tt.wantErr {
					t.Errorf("repo.GetList() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got.Pagination, tt.want.Pagination) {
					t.Errorf("repo.GetList() = %v, want %v", got.Pagination, tt.want.Pagination)
				}

				for i := 0; i < len(got.Data); i++ {
					got.Data[i].ID = 0
				}

				if !reflect.DeepEqual(len(got.Data), len(tt.want.Data)) {
					t.Errorf("repo.GetList() = %v, want %v", len(got.Data), len(tt.want.Data))
				}

				for i := 0; i < len(got.Data); i++ {
					if !reflect.DeepEqual(got.Data[i], tt.want.Data[i]) {
						t.Errorf("repo.GetList() = %v, want %v at index", got.Data[i], tt.want.Data[i])
					}
				}
			})
		}
	})
}
