package product_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/productapi"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func productQueryOk(sd dbtest.SeedData) []dbtest.AppTable {
	prds := make([]product.Product, 0, len(sd.Admins[0].Products)+len(sd.Users[0].Products))
	prds = append(prds, sd.Admins[0].Products...)
	prds = append(prds, sd.Users[0].Products...)

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].ID.String() <= prds[j].ID.String()
	})

	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Admins[0].Token,
			ExpResp: page.Document[productapi.AppProduct]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(prds),
				Items:       toAppProducts(prds),
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
				return cmp.Diff(got, exp)
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
