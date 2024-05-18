package vproductapp

import (
	"time"

	"github.com/ardanlabs/encore/business/domain/vproductbus"
)

// QueryParams represents the set of possible query strings.
type QueryParams struct {
	Page     string `query:"page"`
	Rows     string `query:"rows"`
	OrderBy  string `query:"orderBy"`
	ID       string `query:"product_id"`
	Name     string `query:"name"`
	Cost     string `query:"cost"`
	Quantity string `query:"quantity"`
	UserName string `query:"user_name"`
}

// Product represents information about an individual product with
// extended information.
type Product struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userID"`
	Name        string  `json:"name"`
	Cost        float64 `json:"cost"`
	Quantity    int     `json:"quantity"`
	DateCreated string  `json:"dateCreated"`
	DateUpdated string  `json:"dateUpdated"`
	UserName    string  `json:"userName"`
}

func toAppProduct(prd vproductbus.Product) Product {
	return Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost,
		Quantity:    prd.Quantity,
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
		UserName:    prd.UserName.String(),
	}
}

func toAppProducts(prds []vproductbus.Product) []Product {
	items := make([]Product, len(prds))
	for i, prd := range prds {
		items[i] = toAppProduct(prd)
	}

	return items
}
