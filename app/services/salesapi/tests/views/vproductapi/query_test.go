package vproduct_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/views/vproductapi"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func vproductQueryOk(sd dbtest.SeedData) []dbtest.AppTable {
	prds := toAppVProducts(sd.Admins[0].User, sd.Admins[0].Products)
	prds = append(prds, toAppVProducts(sd.Users[0].User, sd.Users[0].Products)...)

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].ID <= prds[j].ID
	})

	table := []dbtest.AppTable{
		{
			Name:  "all",
			Token: sd.Admins[0].Token,
			ExpResp: page.Document[vproductapi.AppProduct]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(prds),
				Items:       prds,
			},
			ExcFunc: func(ctx context.Context) any {
				qp := vproductapi.QueryParams{
					Page:    1,
					Rows:    10,
					OrderBy: "product_id,ASC",
					Name:    "Name",
				}

				resp, err := salesapi.VProductQuery(ctx, qp)
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
