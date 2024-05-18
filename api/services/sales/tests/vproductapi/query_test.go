package vproduct_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/api/services/sales"
	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/domain/vproductapp"
	"github.com/ardanlabs/encore/app/sdk/query"
	"github.com/google/go-cmp/cmp"
)

func queryOk(sd apitest.SeedData) []apitest.Table {
	prds := toAppVProducts(sd.Admins[0].User, sd.Admins[0].Products)
	prds = append(prds, toAppVProducts(sd.Users[0].User, sd.Users[0].Products)...)

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].ID <= prds[j].ID
	})

	table := []apitest.Table{
		{
			Name:  "all",
			Token: sd.Admins[0].Token,
			ExpResp: query.Result[vproductapp.Product]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(prds),
				Items:       prds,
			},
			ExcFunc: func(ctx context.Context) any {
				qp := vproductapp.QueryParams{
					Page:    "1",
					Rows:    "10",
					OrderBy: "product_id,ASC",
					Name:    "Name",
				}

				resp, err := sales.VProductQuery(ctx, qp)
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
