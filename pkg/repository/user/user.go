package user

import (
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/repository/base"
	"github.com/dwarvesf/go-api/pkg/repository/db"
	"github.com/dwarvesf/go-api/pkg/repository/orm"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type repo struct {
}

func (r *repo) GetList(ctx db.Context, q model.ListQuery) (*model.ListResult[model.User], error) {
	fnSet := base.GetListFuncSet[orm.User, model.User]{
		PrepareQueryFn: func(ctx db.Context, q model.ListQuery) []qm.QueryMod {
			queryParams := []qm.QueryMod{}
			if q.Query != "" {
				queryParams = append(queryParams, qm.Where("lower(name) LIKE lower(?)", "%"+q.Query+"%"))
			}
			return queryParams
		},
		CounableFn: func(q []qm.QueryMod) base.Counable {
			return orm.Users(q...)
		},
		QueryListFn: func(q []qm.QueryMod) ([]*orm.User, error) {
			return orm.Users(q...).All(ctx.Context, ctx.DB)
		},
		MappingFn: toUserModel,
	}
	return base.GetList(ctx, q, fnSet)
}

func (r *repo) Count(ctx db.Context) (int64, error) {
	return orm.Users().Count(ctx.Context, ctx.DB)
}

func (r *repo) GetByID(ctx db.Context, uID int) (*model.User, error) {
	dt, err := orm.FindUser(ctx, ctx.DB, uID)
	return toUserModel(dt), err
}

func (r *repo) GetByEmail(ctx db.Context, email string) (*model.User, error) {
	u, err := orm.Users(
		orm.UserWhere.Email.EQ(email),
	).One(ctx.Context, ctx.DB)
	return toUserModel(u), err
}

func (r *repo) Create(ctx db.Context, user model.SignupRequest) (*model.User, error) {
	u := &orm.User{
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Salt:           user.Salt,
		Status:         string(user.Status),
		Role:           string(user.Role),
		Avatar:         user.Avatar,
	}

	err := u.Insert(ctx, ctx.DB, boil.Infer())
	return toUserModel(u), err
}

func (r *repo) Update(ctx db.Context, uID int, user model.UpdateUserRequest) (*model.User, error) {
	u, err := orm.FindUser(ctx, ctx.DB, uID)
	if err != nil {
		return nil, err
	}

	u.Name = user.FullName
	u.Avatar = user.Avatar

	_, err = u.Update(ctx, ctx.DB, boil.Infer())

	return toUserModel(u), err
}

func (r *repo) UpdatePassword(ctx db.Context, uID int, newPassword string) error {
	u, err := orm.FindUser(ctx, ctx.DB, uID)
	if err != nil {
		return err
	}
	u.HashedPassword = newPassword
	_, err = u.Update(ctx.Context, ctx.DB, boil.Infer())
	return err
}
