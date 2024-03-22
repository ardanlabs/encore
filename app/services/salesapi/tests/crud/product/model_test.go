package product_test

import (
	"time"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/crud/productgrp"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

type tableData struct {
	name       string
	url        string
	token      string
	method     string
	statusCode int
	model      any
	resp       any
	expResp    any
	excFunc    func()
	cmpFunc    func(x interface{}, y interface{}) string
}

type testUser struct {
	user.User
	token    string
	products []product.Product
}

type seedData struct {
	users  []testUser
	admins []testUser
}

func toPointer(r middleware.Response) *middleware.Response {
	return &r
}

func toAppProduct(prd product.Product) productgrp.AppProduct {
	return productgrp.AppProduct{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name,
		Cost:        prd.Cost,
		Quantity:    prd.Quantity,
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProductPtr(prd product.Product) *productgrp.AppProduct {
	appPrd := toAppProduct(prd)
	return &appPrd
}

func toAppProducts(prds []product.Product) []productgrp.AppProduct {
	items := make([]productgrp.AppProduct, len(prds))
	for i, prd := range prds {
		items[i] = toAppProduct(prd)
	}

	return items
}
