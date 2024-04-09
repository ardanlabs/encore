package vproduct_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/apis/services/sales"
	"github.com/ardanlabs/encore/app/api/apptest"
	"github.com/ardanlabs/encore/app/api/page"
	"github.com/ardanlabs/encore/app/core/views/vproductapp"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func vproductQueryOk(sd dbtest.SeedData) []apptest.AppTable {
	prds := toAppVProducts(sd.Admins[0].User, sd.Admins[0].Products)
	prds = append(prds, toAppVProducts(sd.Users[0].User, sd.Users[0].Products)...)

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].ID <= prds[j].ID
	})

	table := []apptest.AppTable{
		{
			Name:  "all",
			Token: sd.Admins[0].Token,
			ExpResp: page.Document[vproductapp.Product]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(prds),
				Items:       prds,
			},
			ExcFunc: func(ctx context.Context) any {
				qp := vproductapp.QueryParams{
					Page:    1,
					Rows:    10,
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
