// Package v1 provides types and support related to web v1 functionality.
package v1

import (
	"errors"

	"encore.dev/beta/errs"
	"encore.dev/middleware"
)

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Err    string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// ErrDetails implements the encore ErrDetails interface.
func (er ErrorResponse) ErrDetails() {}

// =============================================================================

// NewError constructs an encore error based on an app error.
func NewError(code errs.ErrCode, err string) *errs.Error {
	return &errs.Error{
		Code:    code,
		Message: err,
		Details: ErrorResponse{
			Err: err,
		},
	}
}

// NewErrorResponse constructs an encore middleware response with
// a Go error.
func NewErrorResponse(status int, err error) middleware.Response {
	return middleware.Response{
		HTTPStatus: status,
		Err: &errs.Error{
			Code:    errs.Internal,
			Message: err.Error(),
			Details: ErrorResponse{
				Err: err.Error(),
			},
		},
	}
}

// NewErrorResponseWithMessage constructs an encore middleware response
// with a message.
func NewErrorResponseWithMessage(status int, message string) middleware.Response {
	return middleware.Response{
		HTTPStatus: status,
		Err: &errs.Error{
			Code:    errs.Internal,
			Message: message,
			Details: ErrorResponse{
				Err: message,
			},
		},
	}
}

// NewErrorResponseWithFields constructs an encore middleware response
// with an error and fields.
func NewErrorResponseWithFields(status int, message string, fields map[string]string) middleware.Response {
	return middleware.Response{
		HTTPStatus: status,
		Err: &errs.Error{
			Code:    errs.Internal,
			Message: message,
			Details: ErrorResponse{
				Err:    message,
				Fields: fields,
			},
		},
	}
}

// =============================================================================

// TrustedError is used to pass an error during the request through the
// application with web specific context.
type TrustedError struct {
	Err    error
	Status int
}

// NewTrustedError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewTrustedError(err error, status int) error {
	return &TrustedError{err, status}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (te *TrustedError) Error() string {
	return te.Err.Error()
}

// IsTrustedError checks if an error of type TrustedError exists.
func IsTrustedError(err error) bool {
	var te *TrustedError
	return errors.As(err, &te)
}

// GetTrustedError returns a copy of the TrustedError pointer.
func GetTrustedError(err error) *TrustedError {
	var te *TrustedError
	if !errors.As(err, &te) {
		return nil
	}

	return te
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
