package util

import (
	"database/sql"
	"net/http/httptest"
	"testing"

	"github.com/dwarvesf/go-api/pkg/handler/testutil"
	"github.com/dwarvesf/go-api/pkg/handler/v1/viewmodel"
	"github.com/dwarvesf/go-api/pkg/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	type args struct {
		err error
	}
	type expected struct {
		Status int
		Body   string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "valid error",
			args: args{
				err: model.Error{Status: 400, Code: "bad_request", Message: "bad request"},
			},
			expected: expected{
				Status: 400,
				Body:   "bad request",
			},
		},
		{
			name: "valid pointer error",
			args: args{
				err: sql.ErrNoRows,
			},
			expected: expected{

				Status: 500,
				Body:   "no rows",
			},
		},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		ginCtx := testutil.NewRequest(w, testutil.MethodGet, nil, nil, nil, nil)
		t.Run(tt.name, func(t *testing.T) {
			HandleError(ginCtx, tt.args.err)

			assert.Equal(t, tt.expected.Status, w.Code)
			resBody := w.Body.String()
			assert.Contains(t, resBody, tt.expected.Body)
		})
	}
}

func Test_tryParseError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want viewmodel.ErrorResponse
	}{
		{
			name: "valid error",
			args: args{
				err: model.Error{Status: 400, Code: "WRONG_CREDENTIALS", Message: "Wrong username or password"},
			},
			want: viewmodel.ErrorResponse{
				Status:  400,
				Code:    "WRONG_CREDENTIALS",
				Message: "Wrong username or password",
			},
		},
		{
			name: "valid stack error",
			args: args{
				err: errors.WithStack(model.NewError(400, "bad_request", "bad request")),
			},
			want: viewmodel.ErrorResponse{
				Status:  400,
				Code:    "bad_request",
				Message: "bad request",
			},
		},
		{
			name: "valid viewmodel error",
			args: args{
				err: viewmodel.ErrorResponse{
					Status:  400,
					Code:    "bad_request",
					Message: "bad request",
				},
			},
			want: viewmodel.ErrorResponse{
				Status:  400,
				Code:    "bad_request",
				Message: "bad request",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tryParseError(tt.args.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
