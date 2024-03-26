package product_test

import (
	"context"

	"encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/productapi"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func productQueryOk(sd dbtest.SeedData) []dbtest.AppTable {
	total := len(sd.Admins[0].Products) + len(sd.Users[0].Products)

	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Admins[0].Token,
			ExpResp: page.Document[productapi.AppProduct]{
				Page:        1,
				RowsPerPage: 10,
				Total:       total,
				Items:       toAppProducts(append(sd.Admins[0].Products, sd.Users[0].Products...)),
			},
			ExcFunc: func(ctx context.Context) any {
				qp := productapi.QueryParams{
					Page:    1,
					Rows:    10,
					OrderBy: "product_id,ASC",
					Name:    "Name",
				}

				resp, err := salesapi.ProductQuery(ctx, qp)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				if errs, exists := got.(*errs.Error); exists {
					return errs.Message
				}

				gotResp := got.(page.Document[productapi.AppProduct])
				expResp := exp.(page.Document[productapi.AppProduct])

				var found int
				for _, r := range gotResp.Items {
					for _, e := range expResp.Items {
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

func productQueryByIDOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "basic",
			Token:   sd.Users[0].Token,
			ExpResp: toAppProduct(sd.Users[0].Products[0]),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.ProductQueryByID(ctx, sd.Users[0].Products[0].ID.String())
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
