package home_test

import (
	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/crud/homegrp"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/web/page"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func homeQuery200(sd dbtest.SeedData) []dbtest.AppTable {
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
			Name: "basic",
			//url:        "/v1/homes?page=1&rows=10&orderBy=user_id,DESC",
			Token: sd.Admins[0].Token,
			//statusCode: http.StatusOK,
			//method:     http.MethodGet,
			//resp: &page.Document[homegrp.AppHome]{},
			ExpResp: &page.Document[homegrp.AppHome]{
				Page:        1,
				RowsPerPage: 10,
				Total:       total,
				Items:       toAppHomes(append(sd.Admins[0].Homes, sd.Users[0].Homes...)),
			},
			CmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*page.Document[homegrp.AppHome])
				exp := y.(*page.Document[homegrp.AppHome])

				var found int
				for _, r := range resp.Items {
					for _, e := range exp.Items {
						if e.ID == r.ID {
							found++
							break
						}
					}
				}

				if found != total {
					return "number of expected homes didn't match"
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
			//resp:    &homegrp.AppHome{},
			ExpResp: toAppHomePtr(sd.Users[0].Homes[0]),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
