package homeapp

import (
	"errors"

	"github.com/ardanlabs/encore/sales/business/api/order"
	"github.com/ardanlabs/encore/sales/business/core/crud/homebus"
	"github.com/ardanlabs/encore/sales/foundation/validate"
)

func parseOrder(qp QueryParams) (order.By, error) {
	const (
		orderByID     = "home_id"
		orderByType   = "type"
		orderByUserID = "user_id"
	)

	var orderByFields = map[string]string{
		orderByID:     homebus.OrderByID,
		orderByType:   homebus.OrderByType,
		orderByUserID: homebus.OrderByUserID,
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
