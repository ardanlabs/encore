package vproductgrp

import (
	"errors"

	"github.com/ardanlabs/encore/foundation/validate"
)

var errNotProvided = errors.New("not provided")

func validatePaging(qp QueryParams) error {
	if qp.Page <= 0 {
		return validate.NewFieldsError("page", errNotProvided)
	}

	if qp.Rows <= 0 {
		return validate.NewFieldsError("rows", errNotProvided)
	}

	return nil
}
