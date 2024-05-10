package product_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/api/services/sales"
	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/domain/productapp"
	"github.com/ardanlabs/encore/app/sdk/page"
	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/google/go-cmp/cmp"
)

func queryOk(sd apitest.SeedData) []apitest.Table {
	prds := make([]productbus.Product, 0, len(sd.Admins[0].Products)+len(sd.Users[0].Products))
	prds = append(prds, sd.Admins[0].Products...)
	prds = append(prds, sd.Users[0].Products...)

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].ID.String() <= prds[j].ID.String()
	})

	table := []apitest.Table{
		{
			Name:  "all",
			Token: sd.Admins[0].Token,
			ExpResp: page.Document[productapp.Product]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(prds),
				Items:       toAppProducts(prds),
			},
			ExcFunc: func(ctx context.Context) any {
				qp := productapp.QueryParams{
					Page:    "1",
					Rows:    "10",
					OrderBy: "product_id,ASC",
					Name:    "Name",
				}

				resp, err := sales.ProductQuery(ctx, qp)
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

func queryByIDOk(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:    "byid",
			Token:   sd.Users[0].Token,
			ExpResp: toAppProduct(sd.Users[0].Products[0]),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.ProductQueryByID(ctx, sd.Users[0].Products[0].ID.String())
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
