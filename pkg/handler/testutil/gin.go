package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	mw "github.com/dwarvesf/go-api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// HTTPMethod enum
type HTTPMethod string

const (
	// MethodGet http method
	MethodGet HTTPMethod = "GET"
	// MethodPost http method
	MethodPost HTTPMethod = "POST"
	// MethodPut http method
	MethodPut HTTPMethod = "PUT"
	// MethodDelete http method
	MethodDelete HTTPMethod = "Delete"
)

var defaultHeaders = map[string]string{
	"Content-Type": "application/json",
}

// GinContext init a gin context for testing
func GinContext(w http.ResponseWriter) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

func updateHeaders(ctx *gin.Context, headers map[string]string) {
	for k, v := range defaultHeaders {
		ctx.Request.Header.Add(k, v)
	}

	for k, v := range headers {
		ctx.Request.Header.Add(k, v)
	}
}

// NewRequest make a gin.Context request
func NewRequest(w http.ResponseWriter, method HTTPMethod, headers map[string]string, params []gin.Param, u url.Values, body interface{}) *gin.Context {
	ctx := GinContext(w)

	ctx.Request.Method = string(method)
	updateHeaders(ctx, headers)

	if params != nil {
		// set params
		for k := range params {
			v := params[k]
			ctx.Set(v.Key, v.Value)
		}

		ctx.Params = params
	}

	if u != nil {
		// set query params
		ctx.Request.URL.RawQuery = u.Encode()
	}

	if body != nil {
		jsonbytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	}

	return ctx
}

// UpdateJWT update the jwt token
func UpdateJWT(ginCtx *gin.Context, userID int, role string) {

	ctx := context.WithValue(ginCtx.Request.Context(), mw.UserIDCtxKey, userID)
	ctx = context.WithValue(ctx, mw.RoleCtxKey, role)
	ginCtx.Request = ginCtx.Request.WithContext(ctx)
}
