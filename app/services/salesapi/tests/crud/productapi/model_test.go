package product_test

import (
	"time"

	"github.com/ardanlabs/encore/app/services/salesapi/core/crud/productapp"
	"github.com/ardanlabs/encore/business/core/crud/product"
)

func toAppProduct(prd product.Product) productapp.AppProduct {
	return productapp.AppProduct{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name,
		Cost:        prd.Cost,
		Quantity:    prd.Quantity,
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProducts(prds []product.Product) []productapp.AppProduct {
	items := make([]productapp.AppProduct, len(prds))
	for i, prd := range prds {
		items[i] = toAppProduct(prd)
	}

	return items
}
