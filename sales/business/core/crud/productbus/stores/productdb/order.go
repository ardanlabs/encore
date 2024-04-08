package productdb

import (
	"fmt"

	"github.com/ardanlabs/encore/sales/business/api/order"
	"github.com/ardanlabs/encore/sales/business/core/crud/productbus"
)

var orderByFields = map[string]string{
	productbus.OrderByProductID: "product_id",
	productbus.OrderByUserID:    "user_id",
	productbus.OrderByName:      "name",
	productbus.OrderByCost:      "cost",
	productbus.OrderByQuantity:  "quantity",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
