package middleware

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

const subKey = "sub"
const roleKey = "role"

// AuthMiddleware middleware struct for auth
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

// UserIDFromJWTClaims get userID from context
func UserIDFromJWTClaims(jwtClaims map[string]any) (int, error) {
	userID := jwtClaims[subKey]
	if userID == nil {
		return 0, model.ErrInvalidToken
	}
	val, ok := userID.(float64)
	if !ok {
		return 0, model.ErrInvalidToken
	}

	return int(val), nil
}

// WithAuth a middleware to check the access token
func (amw AuthMiddleware) WithAuth(c *gin.Context) {
	jwtClaims, err := amw.authenticate(c)
	if err != nil {
		c.AbortWithStatusJSON(401, err)
		return
	}
	populatedCtx, err := populateContext(c.Request.Context(), jwtClaims)
	c.Request = c.Request.WithContext(populatedCtx)

	c.Next()
}

func populateContext(ctx context.Context, jwtClaims map[string]any) (context.Context, error) {
	ID, ok := jwtClaims[subKey]
	if !ok {
		return ctx, model.ErrInvalidToken
	}
	IDVal, ok := ID.(float64)
	if !ok {
		return ctx, model.ErrInvalidToken
	}

	role, ok := jwtClaims[roleKey]
	if !ok {
		return ctx, model.ErrInvalidToken
	}

	ctx = context.WithValue(ctx, UserIDCtxKey, int(IDVal))
	ctx = context.WithValue(ctx, RoleCtxKey, role)
	return ctx, nil
}

// Authenticate authenticate the request
func (amw AuthMiddleware) Authenticate(c *gin.Context) (map[string]any, error) {
	authHeaderStr := c.Request.Header.Get("Authorization")
	if authHeaderStr == "" {
		return nil, model.ErrNoAuthHeader
	}
	return amw.authenticate(c)
}

func (amw AuthMiddleware) authenticate(c *gin.Context) (map[string]any, error) {
	headers := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(headers) != 2 {
		return nil, model.ErrUnexpectedAuthorizationHeader
	}
	switch headers[0] {
	case "Bearer":
		dt, err := amw.jwtH.ValidateToken(headers[1])
		if err != nil {
			return nil, model.ErrInvalidToken
		}

		return dt, nil
	default:
		return nil, model.ErrUnexpectedAuthorizationHeader
	}
}
