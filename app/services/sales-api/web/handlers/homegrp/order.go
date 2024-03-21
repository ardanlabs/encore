package homegrp

import (
	"errors"

	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/web/order"
	"github.com/ardanlabs/encore/foundation/validate"
)

func parseOrder(qp QueryParams) (order.By, error) {
	const (
		orderByID     = "home_id"
		orderByType   = "type"
		orderByUserID = "user_id"
	)

	var orderByFields = map[string]string{
		orderByID:     home.OrderByID,
		orderByType:   home.OrderByType,
		orderByUserID: home.OrderByUserID,
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
