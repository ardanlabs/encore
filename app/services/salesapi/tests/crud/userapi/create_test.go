package user_test

import (
	"context"
	"net/http"

	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/userapi"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func userCreate200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Admins[0].Token,
			ExpResp: &userapi.AppUser{
				Name:       "Bill Kennedy",
				Email:      "bill@ardanlabs.com",
				Roles:      []string{"ADMIN"},
				Department: "IT",
				Enabled:    true,
			},
			ExcFunc: func(ctx context.Context) any {
				app := userapi.AppNewUser{
					Name:            "Bill Kennedy",
					Email:           "bill@ardanlabs.com",
					Roles:           []string{"ADMIN"},
					Department:      "IT",
					Password:        "123",
					PasswordConfirm: "123",
				}

				resp, err := salesapi.UserCreate(ctx, app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*userapi.AppUser)
				expResp := y.(*userapi.AppUser)

				if _, err := uuid.Parse(resp.ID); err != nil {
					return "bad uuid for ID"
				}

				if resp.DateCreated == "" {
					return "missing date created"
				}

				if resp.DateUpdated == "" {
					return "missing date updated"
				}

				expResp.ID = resp.ID
				expResp.DateCreated = resp.DateCreated
				expResp.DateUpdated = resp.DateUpdated

				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func userCreate400(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "missing-input",
			//url:        "/v1/users",
			Token: sd.Admins[0].Token,
			//method:     http.MethodPost,
			//statusCode: http.StatusBadRequest,
			//model: &userapi.AppNewUser{},
			ExpResp: dbtest.ToPointer(errs.NewResponse(http.StatusBadRequest, validate.FieldErrors{
				validate.FieldError{Field: "email", Err: "email is a required field"},
				validate.FieldError{Field: "name", Err: "name is a required field"},
				validate.FieldError{Field: "password", Err: "password is a required field"},
				validate.FieldError{Field: "roles", Err: "roles is a required field"},
			})),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "bad-role",
			//url:        "/v1/users",
			Token: sd.Admins[0].Token,
			//method:     http.MethodPost,
			//statusCode: http.StatusBadRequest,
			// model: &userapi.AppNewUser{
			// 	Name:            "Bill Kennedy",
			// 	Email:           "bill@ardanlabs.com",
			// 	Roles:           []string{"BAD ROLE"},
			// 	Department:      "IT",
			// 	Password:        "123",
			// 	PasswordConfirm: "123",
			// },
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusBadRequest, `parse: invalid role \"BAD ROLE\"`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func userCreate401(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "emptytoken",
			//url:        "/v1/users",
			Token: "",
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "badtoken",
			//url:        "/v1/users",
			Token: sd.Admins[0].Token[:10],
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "badsig",
			//url:        "/v1/users",
			Token: sd.Admins[0].Token + "A",
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "wronguser",
			//url:        "/v1/users",
			Token: sd.Users[0].Token,
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
