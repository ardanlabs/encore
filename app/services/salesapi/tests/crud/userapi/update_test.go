package user_test

import (
	"net/http"

	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/userapi"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func userUpdate200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "basic",
			//url:        fmt.Sprintf("/v1/users/%s", sd.users[0].ID),
			Token: sd.Users[0].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusOK,
			// model: &userapi.AppUpdateUser{
			// 	Name:            dbtest.StringPointer("Jack Kennedy"),
			// 	Email:           dbtest.StringPointer("jack@ardanlabs.com"),
			// 	Roles:           []string{"ADMIN"},
			// 	Department:      dbtest.StringPointer("IT"),
			// 	Password:        dbtest.StringPointer("123"),
			// 	PasswordConfirm: dbtest.StringPointer("123"),
			// },
			//resp: &userapi.AppUser{},
			ExpResp: &userapi.AppUser{
				Name:       "Jack Kennedy",
				Email:      "jack@ardanlabs.com",
				Roles:      []string{"ADMIN"},
				Department: "IT",
				Enabled:    true,
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

func userUpdate400(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "bad-input",
			//url:        fmt.Sprintf("/v1/users/%s", sd.users[0].ID),
			Token: sd.Users[0].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusBadRequest,
			// model: &userapi.AppUpdateUser{
			// 	Email:           dbtest.StringPointer("bill@"),
			// 	PasswordConfirm: dbtest.StringPointer("jack"),
			// },
			//resp: &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponse(http.StatusBadRequest, validate.FieldErrors{
				validate.FieldError{Field: "email", Err: "email must be a valid email address"},
				validate.FieldError{Field: "passwordConfirm", Err: "passwordConfirm must be equal to Password"},
			})),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "bad-role",
			//url:        fmt.Sprintf("/v1/users/%s", sd.users[0].ID),
			Token: sd.Users[0].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusBadRequest,
			// model: &userapi.AppUpdateUser{
			// 	Roles: []string{"BAD ROLE"},
			// },
			//resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusBadRequest, `parse: invalid role \"BAD ROLE\"`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func userUpdate401(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "emptytoken",
			//url:        fmt.Sprintf("/v1/users/%s", sd.users[0].ID),
			Token: "",
			//method:     http.MethodPut,
			////statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "badsig",
			//url:        fmt.Sprintf("/v1/users/%s", sd.users[0].ID),
			Token: sd.Users[0].Token + "A",
			//method:     http.MethodPut,
			//statusCode: http.StatusUnauthorized,
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "wronguser",
			//url:        fmt.Sprintf("/v1/users/%s", sd.admins[0].ID),
			Token: sd.Users[0].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusUnauthorized,
			// model: &userapi.AppUpdateUser{
			// 	Name:            dbtest.StringPointer("Bill Kennedy"),
			// 	Email:           dbtest.StringPointer("bill@ardanlabs.com"),
			// 	Roles:           []string{"ADMIN"},
			// 	Department:      dbtest.StringPointer("IT"),
			// 	Password:        dbtest.StringPointer("123"),
			// 	PasswordConfirm: dbtest.StringPointer("123"),
			// },
			//resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
