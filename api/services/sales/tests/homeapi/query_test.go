package home_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/api/services/sales"
	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/domain/homeapp"
	"github.com/ardanlabs/encore/app/sdk/query"
	"github.com/ardanlabs/encore/business/domain/homebus"
	"github.com/google/go-cmp/cmp"
)

func queryOk(sd apitest.SeedData) []apitest.Table {
	hmes := make([]homebus.Home, 0, len(sd.Admins[0].Homes)+len(sd.Users[0].Homes))
	hmes = append(hmes, sd.Admins[0].Homes...)
	hmes = append(hmes, sd.Users[0].Homes...)

	sort.Slice(hmes, func(i, j int) bool {
		return hmes[i].ID.String() <= hmes[j].ID.String()
	})

	table := []apitest.Table{
		{
			Name:  "all",
			Token: sd.Admins[0].Token,
			ExpResp: query.Result[homeapp.Home]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(hmes),
				Items:       toAppHomes(hmes),
			},
			ExcFunc: func(ctx context.Context) any {
				qp := homeapp.QueryParams{
					Page:    "1",
					Rows:    "10",
					OrderBy: "home_id,ASC",
				}

				resp, err := sales.HomeQuery(ctx, qp)
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
			ExpResp: toAppHome(sd.Users[0].Homes[0]),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.HomeQueryByID(ctx, sd.Users[0].Homes[0].ID.String())
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
