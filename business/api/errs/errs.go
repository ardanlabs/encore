// Package errs provides types and support related to web error functionality.
package errs

import (
	"fmt"
	"net/http"

	"encore.dev/beta/errs"
	"encore.dev/middleware"
)

// ExtraDetails provides the caller with more error context.
type ExtraDetails struct {
	HTTPStatusCode int    `json:"httpStatusCode"`
	HTTPStatus     string `json:"httpStatus"`
}

func (ExtraDetails) ErrDetails() {}

// New constructs an encore error based on an app error.
func New(code errs.ErrCode, err error) *errs.Error {
	return &errs.Error{
		Code:    code,
		Message: err.Error(),
		Details: ExtraDetails{
			HTTPStatusCode: code.HTTPStatus(),
			HTTPStatus:     http.StatusText(code.HTTPStatus()),
		},
	}
}

// Newf constructs an encore error based on a error message.
func Newf(code errs.ErrCode, format string, v ...any) *errs.Error {
	return &errs.Error{
		Code:    code,
		Message: fmt.Sprintf(format, v...),
		Details: ExtraDetails{
			HTTPStatusCode: code.HTTPStatus(),
			HTTPStatus:     http.StatusText(code.HTTPStatus()),
		},
	}
}

// NewResponse constructs an encore middleware response with a Go error.
func NewResponse(code errs.ErrCode, err error) middleware.Response {
	return middleware.Response{
		Err: &errs.Error{
			Code:    code,
			Message: err.Error(),
			Details: ExtraDetails{
				HTTPStatusCode: code.HTTPStatus(),
				HTTPStatus:     http.StatusText(code.HTTPStatus()),
			},
		},
	}
}

// NewResponsef constructs an encore middleware response with a message.
func NewResponsef(code errs.ErrCode, format string, v ...any) middleware.Response {
	return middleware.Response{
		Err: &errs.Error{
			Code:    code,
			Message: fmt.Sprintf(format, v...),
			Details: ExtraDetails{
				HTTPStatusCode: code.HTTPStatus(),
				HTTPStatus:     http.StatusText(code.HTTPStatus()),
			},
		},
	}
}
