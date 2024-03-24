package product_test

import (
	"time"

	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/productapi"
	"github.com/ardanlabs/encore/business/core/crud/product"
)

func toAppProduct(prd product.Product) productapi.AppProduct {
	return productapi.AppProduct{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name,
		Cost:        prd.Cost,
		Quantity:    prd.Quantity,
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProductPtr(prd product.Product) *productapi.AppProduct {
	appPrd := toAppProduct(prd)
	return &appPrd
}

func toAppProducts(prds []product.Product) []productapi.AppProduct {
	items := make([]productapi.AppProduct, len(prds))
	for i, prd := range prds {
		items[i] = toAppProduct(prd)
	}

	return items
}
