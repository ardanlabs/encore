package user_test

import (
	"context"

	"encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/userapi"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func userQuery200(sd dbtest.SeedData) []dbtest.AppTable {
	usrs := make([]user.User, 0, len(sd.Admins)+len(sd.Users))
	usrsMap := make(map[uuid.UUID]user.User)

	for _, adm := range sd.Admins {
		usrsMap[adm.ID] = adm.User
		usrs = append(usrs, adm.User)
	}

	for _, usr := range sd.Users {
		usrsMap[usr.ID] = usr.User
		usrs = append(usrs, usr.User)
	}

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
				if errs, exists := got.(*errs.Error); exists {
					return errs.Message
				}

				gotResp := got.(page.Document[userapi.AppUser])
				expResp := exp.(page.Document[userapi.AppUser])

				var found int
				for _, r := range gotResp.Items {
					for _, e := range expResp.Items {
						if e.ID == r.ID {
							found++
							break
						}
					}
				}

				if found != len(usrs) {
					return "number of expected users didn't match"
				}

				return ""
			},
		},
	}

	return table
}

func userQueryByID200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "basic",
			Token:   sd.Users[0].Token,
			ExpResp: toAppUserPtr(sd.Users[0].User),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.UserQueryByID(ctx, sd.Users[0].ID.String())
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(x any, y any) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
