package vproduct_test

import (
	"github.com/ardanlabs/encore/app/services/salesapi/apis/views/vproductapi"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/data/dbtest"
)

func vproductQuery200(sd dbtest.SeedData) []dbtest.AppTable {
	total := len(sd.Admins[0].Products) + len(sd.Users[0].Products)

	allProducts := toAppVProducts(sd.Admins[0].User, sd.Admins[0].Products)
	allProducts = append(allProducts, toAppVProducts(sd.Users[0].User, sd.Users[0].Products)...)

	table := []dbtest.AppTable{
		{
			Name: "basic",
			//url:        "/v1/vproducts?page=1&rows=10&orderBy=product_id,DESC",
			Token: sd.Admins[0].Token,
			//statusCode: http.StatusOK,
			//method:     http.MethodGet,
			//resp:       &page.Document[vproductapi.AppProduct]{},
			ExpResp: &page.Document[vproductapi.AppProduct]{
				Page:        1,
				RowsPerPage: 10,
				Total:       total,
				Items:       allProducts,
			},
			CmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*page.Document[vproductapi.AppProduct])
				exp := y.(*page.Document[vproductapi.AppProduct])

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
