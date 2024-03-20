// Package v1 provides types and support related to web v1 functionality.
package v1

import (
	"net/http"

	"encore.dev/beta/errs"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/foundation/validate"
)

// NewError constructs an encore error based on an app error.
func NewError(status int, err error) *errs.Error {
	return &errs.Error{
		Code:    errMap[status],
		Message: err.Error(),
	}
}

// NewErrorWithMessage constructs an encore error based on a error message.
func NewErrorWithMessage(status int, message string) *errs.Error {
	return &errs.Error{
		Code:    errMap[status],
		Message: message,
	}
}

type extraDetails struct {
	HTTPStatusCode int                  `json:"httpStatusCode"`
	HTTPStatus     string               `json:"httpStatus"`
	Fields         validate.FieldErrors `json:"fields,omitempty"`
}

func (extraDetails) ErrDetails() {}

// NewErrorResponse constructs an encore middleware response with
// a Go error.
func NewErrorResponse(httpStatus int, err error) middleware.Response {
	return middleware.Response{
		HTTPStatus: httpStatus,
		Err: &errs.Error{
			Code:    errMap[httpStatus],
			Message: err.Error(),
			Details: extraDetails{
				HTTPStatusCode: httpStatus,
				HTTPStatus:     http.StatusText(httpStatus),
			},
		},
	}
}

// NewErrorResponseWithMessage constructs an encore middleware response
// with a message.
func NewErrorResponseWithMessage(httpStatus int, message string) middleware.Response {
	return middleware.Response{
		HTTPStatus: httpStatus,
		Err: &errs.Error{
			Code:    errMap[httpStatus],
			Message: message,
			Details: extraDetails{
				HTTPStatusCode: httpStatus,
				HTTPStatus:     http.StatusText(httpStatus),
			},
		},
	}
}

// NewErrorResponseWithFields constructs an encore middleware response
// with an error and fields.
func NewErrorResponseWithFields(httpStatus int, message string, fields validate.FieldErrors) middleware.Response {
	return middleware.Response{
		HTTPStatus: httpStatus,
		Err: &errs.Error{
			Code:    errMap[httpStatus],
			Message: message,
			Details: extraDetails{
				HTTPStatusCode: httpStatus,
				HTTPStatus:     http.StatusText(httpStatus),
				Fields:         fields,
			},
		},
	}
}

// =============================================================================

// PageDocument is the form used for API responses from query API calls.
type PageDocument[T any] struct {
	Items       []T `json:"items"`
	Total       int `json:"total"`
	Page        int `json:"page"`
	RowsPerPage int `json:"rowsPerPage"`
}

// NewPageDocument constructs a response value for a web paging trusted.
func NewPageDocument[T any](items []T, total int, page int, rowsPerPage int) PageDocument[T] {
	return PageDocument[T]{
		Items:       items,
		Total:       total,
		Page:        page,
		RowsPerPage: rowsPerPage,
	}
}

// =============================================================================

var errMap = map[int]errs.ErrCode{
	http.StatusOK:                  errs.OK,
	http.StatusInternalServerError: errs.Internal,
	http.StatusBadRequest:          errs.FailedPrecondition,
	http.StatusGatewayTimeout:      errs.DeadlineExceeded,
	http.StatusNotFound:            errs.NotFound,
	http.StatusConflict:            errs.Aborted,
	http.StatusForbidden:           errs.PermissionDenied,
	http.StatusTooManyRequests:     errs.ResourceExhausted,
	http.StatusNotImplemented:      errs.Unimplemented,
	http.StatusServiceUnavailable:  errs.Unavailable,
	http.StatusUnauthorized:        errs.Unauthenticated,
}
