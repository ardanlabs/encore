package home_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/homeapi"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func homeQueryOk(sd dbtest.SeedData) []dbtest.AppTable {
	hmes := make([]home.Home, 0, len(sd.Admins[0].Homes)+len(sd.Users[0].Homes))
	hmes = append(hmes, sd.Admins[0].Homes...)
	hmes = append(hmes, sd.Users[0].Homes...)

	sort.Slice(hmes, func(i, j int) bool {
		return hmes[i].ID.String() <= hmes[j].ID.String()
	})

	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Admins[0].Token,
			ExpResp: page.Document[homeapi.AppHome]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(hmes),
				Items:       toAppHomes(hmes),
			},
			ExcFunc: func(ctx context.Context) any {
				qp := homeapi.QueryParams{
					Page:    1,
					Rows:    10,
					OrderBy: "home_id,ASC",
				}

				resp, err := salesapi.HomeQuery(ctx, qp)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				if _, exists := got.(page.Document[homeapi.AppHome]); !exists {
					return "error occurred"
				}

				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func homeQueryByIDOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "basic",
			Token:   sd.Users[0].Token,
			ExpResp: toAppHome(sd.Users[0].Homes[0]),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.HomeQueryByID(ctx, sd.Users[0].Homes[0].ID.String())
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				if _, exists := got.(homeapi.AppHome); !exists {
					return "error occurred"
				}

				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
