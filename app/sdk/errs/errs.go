// Package errs provides types and support related to web error functionality.
package errs

import (
	"encoding/json"
	"errors"
	"fmt"

	"encore.dev/beta/errs"
	"encore.dev/middleware"
)

// New constructs an encore error based on an app error.
func New(code errs.ErrCode, err error) *errs.Error {
	return &errs.Error{
		Code:    code,
		Message: err.Error(),
	}
}

// Newf constructs an encore error based on a error message.
func Newf(code errs.ErrCode, format string, v ...any) *errs.Error {
	return &errs.Error{
		Code:    code,
		Message: fmt.Sprintf(format, v...),
	}
}

// NewResponse constructs an encore middleware response with a Go error.
func NewResponse(code errs.ErrCode, err error) middleware.Response {
	return middleware.Response{
		Err: &errs.Error{
			Code:    code,
			Message: err.Error(),
		},
	}
}

// NewResponsef constructs an encore middleware response with a message.
func NewResponsef(code errs.ErrCode, format string, v ...any) middleware.Response {
	return middleware.Response{
		Err: &errs.Error{
			Code:    code,
			Message: fmt.Sprintf(format, v...),
		},
	}
}

// =============================================================================

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

// FieldErrors represents a collection of field errors.
type FieldErrors []FieldError

// NewFieldsError creates an fields error.
func NewFieldsError(field string, err error) error {
	return FieldErrors{
		{
			Field: field,
			Err:   err.Error(),
		},
	}
}

// Error implements the error interface.
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

// Fields returns the fields that failed validation
func (fe FieldErrors) Fields() map[string]string {
	m := make(map[string]string, len(fe))
	for _, fld := range fe {
		m[fld.Field] = fld.Err
	}
	return m
}

// IsFieldErrors checks if an error of type FieldErrors exists.
func IsFieldErrors(err error) bool {
	var fe FieldErrors
	return errors.As(err, &fe)
}

// GetFieldErrors returns a copy of the FieldErrors pointer.
func GetFieldErrors(err error) FieldErrors {
	var fe FieldErrors
	if !errors.As(err, &fe) {
		return nil
	}
	return fe
}
