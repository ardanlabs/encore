package vproductdb

import (
	"fmt"

	"github.com/ardanlabs/encore/business/domain/vproductbus"
	"github.com/ardanlabs/encore/business/sdk/order"
)

var orderByFields = map[string]string{
	vproductbus.OrderByProductID: "product_id",
	vproductbus.OrderByUserID:    "user_id",
	vproductbus.OrderByName:      "name",
	vproductbus.OrderByCost:      "cost",
	vproductbus.OrderByQuantity:  "quantity",
	vproductbus.OrderByUserName:  "user_name",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
