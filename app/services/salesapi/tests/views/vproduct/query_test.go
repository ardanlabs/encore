package vproduct_test

import (
	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/views/vproductgrp"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/web/page"
)

func vproductQuery200(sd dbtest.SeedData) []dbtest.AppTable {
	total := len(sd.Admins[1].Products) + len(sd.Users[1].Products)

	allProducts := toAppVProducts(sd.Admins[1].User, sd.Admins[1].Products)
	allProducts = append(allProducts, toAppVProducts(sd.Users[1].User, sd.Users[1].Products)...)

	table := []dbtest.AppTable{
		{
			Name: "basic",
			//url:        "/v1/vproducts?page=1&rows=10&orderBy=product_id,DESC",
			Token: sd.Admins[1].Token,
			//statusCode: http.StatusOK,
			//method:     http.MethodGet,
			//resp:       &page.Document[vproductgrp.AppProduct]{},
			ExpResp: &page.Document[vproductgrp.AppProduct]{
				Page:        1,
				RowsPerPage: 10,
				Total:       total,
				Items:       allProducts,
			},
			CmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*page.Document[vproductgrp.AppProduct])
				exp := y.(*page.Document[vproductgrp.AppProduct])

				var found int
				for _, r := range resp.Items {
					for _, e := range exp.Items {
						if e.ID == r.ID && e.UserName == r.UserName {
							found++
							break
						}
					}
				}

				if found != total {
					return "number of expected products didn't match"
				}

				return ""
			},
		},
	}

	return table
}
