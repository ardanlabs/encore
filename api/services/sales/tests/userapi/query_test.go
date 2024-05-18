package user_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/api/services/sales"
	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/domain/userapp"
	"github.com/ardanlabs/encore/app/sdk/query"
	"github.com/ardanlabs/encore/business/domain/userbus"
	"github.com/google/go-cmp/cmp"
)

func queryOk(sd apitest.SeedData) []apitest.Table {
	usrs := make([]userbus.User, 0, len(sd.Admins)+len(sd.Users))

	for _, adm := range sd.Admins {
		usrs = append(usrs, adm.User)
	}

	for _, usr := range sd.Users {
		usrs = append(usrs, usr.User)
	}

	sort.Slice(usrs, func(i, j int) bool {
		return usrs[i].ID.String() <= usrs[j].ID.String()
	})

	table := []apitest.Table{
		{
			Name:  "all",
			Token: sd.Admins[0].Token,
			ExpResp: query.Result[userapp.User]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(usrs),
				Items:       toAppUsers(usrs),
			},
			ExcFunc: func(ctx context.Context) any {
				qp := userapp.QueryParams{
					Page:    "1",
					Rows:    "10",
					OrderBy: "user_id,ASC",
					Name:    "Name",
				}

				resp, err := sales.UserQuery(ctx, qp)
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
			ExpResp: toAppUser(sd.Users[0].User),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.UserQueryByID(ctx, sd.Users[0].ID.String())
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
