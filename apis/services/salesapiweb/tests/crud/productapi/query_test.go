package product_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/apis/services/salesapiweb"
	"github.com/ardanlabs/encore/app/api/apptest"
	"github.com/ardanlabs/encore/app/api/page"
	"github.com/ardanlabs/encore/app/core/crud/productapp"
	"github.com/ardanlabs/encore/business/core/crud/productbus"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func productQueryOk(sd dbtest.SeedData) []apptest.AppTable {
	prds := make([]productbus.Product, 0, len(sd.Admins[0].Products)+len(sd.Users[0].Products))
	prds = append(prds, sd.Admins[0].Products...)
	prds = append(prds, sd.Users[0].Products...)

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].ID.String() <= prds[j].ID.String()
	})

	table := []apptest.AppTable{
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
					Page:    1,
					Rows:    10,
					OrderBy: "product_id,ASC",
					Name:    "Name",
				}

				resp, err := salesapiweb.ProductQuery(ctx, qp)
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

func productQueryByIDOk(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "byid",
			Token:   sd.Users[0].Token,
			ExpResp: toAppProduct(sd.Users[0].Products[0]),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapiweb.ProductQueryByID(ctx, sd.Users[0].Products[0].ID.String())
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
