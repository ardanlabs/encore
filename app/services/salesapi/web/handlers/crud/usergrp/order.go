package usergrp

import (
	"errors"

	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/web/order"
	"github.com/ardanlabs/encore/foundation/validate"
)

func parseOrder(qp QueryParams) (order.By, error) {
	const (
		orderByID      = "user_id"
		orderByName    = "name"
		orderByEmail   = "email"
		orderByRoles   = "roles"
		orderByEnabled = "enabled"
	)

	var orderByFields = map[string]string{
		orderByID:      user.OrderByID,
		orderByName:    user.OrderByName,
		orderByEmail:   user.OrderByEmail,
		orderByRoles:   user.OrderByRoles,
		orderByEnabled: user.OrderByEnabled,
	}

	orderBy, err := order.Parse(qp.OrderBy, order.NewBy(orderByID, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
