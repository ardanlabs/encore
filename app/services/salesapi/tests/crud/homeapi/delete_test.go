package home_test

import (
	"net/http"

	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func homeDelete200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "asuser",
			////url:        fmt.Sprintf("/v1/homes/%s", sd.users[0].homes[0].ID),
			Token: sd.Users[0].Token,
			//method:     http.MethodDelete,
			//statusCode: http.StatusNoContent,
		},
		{
			Name: "asadmin",
			////url:        fmt.Sprintf("/v1/homes/%s", sd.admins[0].homes[0].ID),
			Token: sd.Admins[0].Token,
			//method:     http.MethodDelete,
			//statusCode: http.StatusNoContent,
		},
	}

	return table
}

func homeDelete401(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "emptytoken",
			//url:   fmt.Sprintf("/v1/homes/%s", sd.users[0].homes[1].ID),
			Token: "",
			//method:     http.MethodDelete,
			//statusCode: http.StatusUnauthorized,
			//resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "badsig",
			//url:   fmt.Sprintf("/v1/homes/%s", sd.users[0].homes[1].ID),
			Token: sd.Users[0].Token + "A",
			//method:     http.MethodDelete,
			//statusCode: http.StatusUnauthorized,
			//resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "wronguser",
			//url:   fmt.Sprintf("/v1/homes/%s", sd.users[0].homes[1].ID),
			Token: sd.Users[1].Token,
			//method:     http.MethodDelete,
			//statusCode: http.StatusUnauthorized,
			//resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
