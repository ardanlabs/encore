package vproduct_test

import (
	"time"

	"github.com/ardanlabs/encore/app/services/salesapi/core/views/vproductapp"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

func toAppVProduct(usr user.User, prd product.Product) vproductapp.AppProduct {
	return vproductapp.AppProduct{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name,
		Cost:        prd.Cost,
		Quantity:    prd.Quantity,
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
		UserName:    usr.Name,
	}
}

func toAppVProducts(usr user.User, prds []product.Product) []vproductapp.AppProduct {
	items := make([]vproductapp.AppProduct, len(prds))
	for i, prd := range prds {
		items[i] = toAppVProduct(usr, prd)
	}

	return items
}
