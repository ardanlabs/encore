package user_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/core/crud/userapp"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func userQueryOk(sd dbtest.SeedData) []dbtest.AppTable {
	usrs := make([]user.User, 0, len(sd.Admins)+len(sd.Users))

	for _, adm := range sd.Admins {
		usrs = append(usrs, adm.User)
	}

	for _, usr := range sd.Users {
		usrs = append(usrs, usr.User)
	}

	sort.Slice(usrs, func(i, j int) bool {
		return usrs[i].ID.String() <= usrs[j].ID.String()
	})

	table := []dbtest.AppTable{
		{
			Name:  "all",
			Token: sd.Admins[0].Token,
			ExpResp: page.Document[userapp.AppUser]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(usrs),
				Items:       toAppUsers(usrs),
			},
			ExcFunc: func(ctx context.Context) any {
				qp := userapp.QueryParams{
					Page:    1,
					Rows:    10,
					OrderBy: "user_id,ASC",
					Name:    "Name",
				}

				resp, err := salesapi.UserQuery(ctx, qp)
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

func userQueryByIDOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "byid",
			Token:   sd.Users[0].Token,
			ExpResp: toAppUser(sd.Users[0].User),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.UserQueryByID(ctx, sd.Users[0].ID.String())
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
