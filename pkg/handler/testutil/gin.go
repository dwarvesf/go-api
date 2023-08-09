package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
)

// HttpMethod enum
type HttpMethod string

const (
	MethodGet    HttpMethod = "GET"
	MethodPost   HttpMethod = "POST"
	MethodPut    HttpMethod = "PUT"
	MethodDelete HttpMethod = "Delete"
)

var defaultHeaders = map[string]string{
	"Content-Type": "application/json",
}

// GinContext init a gin context for testing
func GinContext(w *httptest.ResponseRecorder) *gin.Context {
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
func NewRequest(w *httptest.ResponseRecorder, method HttpMethod, headers map[string]string, params []gin.Param, u url.Values, body interface{}) *gin.Context {
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
