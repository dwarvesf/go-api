package mw

import (
	"context"
	"strings"

	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/dwarvesf/go-api/pkg/service/jwthelper"
	"github.com/gin-gonic/gin"
)

type contextKey string

// RoleCtxKey is the key used to store to context
const RoleCtxKey = contextKey("role")

// UserIDCtxKey is the key used to store to context
const UserIDCtxKey = contextKey("userID")

type AuthMiddleware struct {
	jwtH jwthelper.Helper
}

// NewAuthMiddleware new middleware
func NewAuthMiddleware(secret string) AuthMiddleware {
	return AuthMiddleware{
		jwtH: jwthelper.NewHelper(secret),
	}
}

// UserIDFromContext get userID from context
func UserIDFromContext(ctx context.Context) (int, error) {
	userID := ctx.Value(UserIDCtxKey)
	if userID == nil {
		return 0, model.ErrInvalidToken
	}
	val, ok := userID.(int)
	if !ok {
		return 0, model.ErrInvalidToken
	}

	return val, nil

}

// WithAuth a middleware to check the access token
func (amw *AuthMiddleware) WithAuth(c *gin.Context) {
	err := amw.authenticate(c)
	if err != nil {
		c.AbortWithStatusJSON(401, err)
		return
	}

	c.Next()
}

func (amw *AuthMiddleware) authenticate(c *gin.Context) error {
	headers := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(headers) != 2 {
		return model.ErrUnexpectedAuthorizationHeader
	}
	switch headers[0] {
	case "Bearer":
		dt, err := amw.jwtH.ValidateToken(headers[1])
		if err != nil {
			return model.ErrInvalidToken
		}
		ID, ok := dt["sub"]
		if !ok {
			return model.ErrInvalidToken
		}
		IDVal, ok := ID.(float64)
		if !ok {
			return model.ErrInvalidToken
		}

		role, ok := dt["role"]
		if !ok {
			return model.ErrInvalidToken
		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, UserIDCtxKey, int(IDVal))
		ctx = context.WithValue(ctx, RoleCtxKey, role)

		c.Request = c.Request.WithContext(ctx)

		return nil
	default:
		return model.ErrUnexpectedAuthorizationHeader
	}
}
