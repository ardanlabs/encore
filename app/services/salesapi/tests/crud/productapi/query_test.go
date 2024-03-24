package product_test

import (
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/productapi"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/web/page"
	"github.com/google/go-cmp/cmp"
)

func productQuery200(sd dbtest.SeedData) []dbtest.AppTable {
	total := len(sd.Admins[1].Products) + len(sd.Users[1].Products)

	table := []dbtest.AppTable{
		{
			Name: "basic",
			//url:        "/v1/products?page=1&rows=10&orderBy=product_id,DESC",
			Token: sd.Admins[1].Token,
			//statusCode: http.StatusOK,
			//method:     http.MethodGet,
			//resp:       &page.Document[productapi.AppProduct]{},
			ExpResp: &page.Document[productapi.AppProduct]{
				Page:        1,
				RowsPerPage: 10,
				Total:       total,
				Items:       toAppProducts(append(sd.Admins[1].Products, sd.Users[1].Products...)),
			},
			CmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*page.Document[productapi.AppProduct])
				exp := y.(*page.Document[productapi.AppProduct])

				var found int
				for _, r := range resp.Items {
					for _, e := range exp.Items {
						if e.ID == r.ID {
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

func productQueryByID200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "basic",
			//url:        fmt.Sprintf("/v1/products/%s", sd.Users[1].products[0].ID),
			Token: sd.Users[1].Token,
			//statusCode: http.StatusOK,
			//method:     http.MethodGet,
			//resp:       &productapi.AppProduct{},
			ExpResp: toAppProductPtr(sd.Users[1].Products[0]),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
