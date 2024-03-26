package home_test

import (
	"context"

	"encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/homeapi"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func homeQueryOk(sd dbtest.SeedData) []dbtest.AppTable {
	total := len(sd.Admins[0].Homes) + len(sd.Users[0].Homes)
	usrsMap := make(map[uuid.UUID]user.User)

	for _, adm := range sd.Admins {
		usrsMap[adm.ID] = adm.User
	}
	for _, usr := range sd.Users {
		usrsMap[usr.ID] = usr.User
	}

	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Admins[0].Token,
			ExpResp: &page.Document[homeapi.AppHome]{
				Page:        1,
				RowsPerPage: 10,
				Total:       total,
				Items:       toAppHomes(append(sd.Admins[0].Homes, sd.Users[0].Homes...)),
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
				if errs, exists := got.(*errs.Error); exists {
					return errs.Message
				}

				gotResp := got.(page.Document[homeapi.AppHome])
				expResp := exp.(page.Document[homeapi.AppHome])

				var found int
				for _, r := range gotResp.Items {
					for _, e := range expResp.Items {
						if e.ID == r.ID {
							found++
							break
						}
					}
				}

				if found != total {
					return "number of expected products didn't match"
				}

				return ""
			},
		},
	}

	return table
}

func homeQueryByID200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "basic",
			//url:        fmt.Sprintf("/v1/homes/%s", sd.Users[0].Homes[0].ID),
			Token: sd.Users[0].Token,
			//statusCode: http.StatusOK,
			//method:     http.MethodGet,
			//resp:    &homeapi.AppHome{},
			ExpResp: toAppHomePtr(sd.Users[0].Homes[0]),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
