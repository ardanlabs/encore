package user_test

import (
	"context"
	"sort"

	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/userapi"
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
		return usrs[i].Name <= usrs[j].Name
	})

	table := []dbtest.AppTable{
		{
			Name:  "query",
			Token: sd.Admins[0].Token,
			ExpResp: page.Document[userapi.AppUser]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(usrs),
				Items:       toAppUsers(usrs),
			},
			ExcFunc: func(ctx context.Context) any {
				qp := userapi.QueryParams{
					Page:    1,
					Rows:    10,
					OrderBy: "name,ASC",
					Name:    "Name",
				}

				resp, err := salesapi.UserQuery(ctx, qp)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(page.Document[userapi.AppUser])
				if !exists {
					return "error occurred"
				}

				expResp := exp.(page.Document[userapi.AppUser])

				var found int
				for i := range gotResp.Items {
					for j := range expResp.Items {
						if expResp.Items[i].ID == gotResp.Items[j].ID {
							found++
						}
					}
				}

				if found != len(usrs) {
					return "number of expected users didn't match"
				}

				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func userQueryByIDOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "basic",
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
				if _, exists := got.(userapi.AppUser); !exists {
					return "error occurred"
				}

				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
