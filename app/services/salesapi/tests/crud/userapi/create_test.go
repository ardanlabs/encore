package user_test

import (
	"context"
	"net/http"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/userapi"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func userCreateOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Admins[0].Token,
			ExpResp: userapi.AppUser{
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
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(userapi.AppUser)
				expResp := exp.(userapi.AppUser)

				if _, err := uuid.Parse(gotResp.ID); err != nil {
					return "bad uuid for ID"
				}

				if gotResp.DateCreated == "" {
					return "missing date created"
				}

				if gotResp.DateUpdated == "" {
					return "missing date updated"
				}

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func userCreateBad(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "missing-input",
			Token:   sd.Admins[0].Token,
			ExpResp: errs.Newf(http.StatusBadRequest, "validate: [{\"field\":\"name\",\"error\":\"name is a required field\"},{\"field\":\"email\",\"error\":\"email is a required field\"},{\"field\":\"roles\",\"error\":\"roles is a required field\"},{\"field\":\"password\",\"error\":\"password is a required field\"}]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.UserCreate(ctx, userapi.AppNewUser{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*eerrs.Error)
				expResp := exp.(*eerrs.Error)

				return dbtest.CmpErrors(gotResp, expResp)
			},
		},
		{
			Name:    "bad-role",
			Token:   sd.Admins[0].Token,
			ExpResp: errs.Newf(http.StatusBadRequest, "parse: invalid role \"BAD ROLE\""),
			ExcFunc: func(ctx context.Context) any {
				app := userapi.AppNewUser{
					Name:            "Bill Kennedy",
					Email:           "bill2@ardanlabs.com",
					Roles:           []string{"BAD ROLE"},
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
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*eerrs.Error)
				expResp := exp.(*eerrs.Error)

				return dbtest.CmpErrors(gotResp, expResp)
			},
		},
	}

	return table
}

func userCreateAuth(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(http.StatusUnauthorized, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.UserCreate(ctx, userapi.AppNewUser{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*eerrs.Error)
				expResp := exp.(*eerrs.Error)

				return dbtest.CmpErrors(gotResp, expResp)
			},
		},
		{
			Name:    "badtoken",
			Token:   sd.Admins[0].Token[:10],
			ExpResp: errs.Newf(http.StatusUnauthorized, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.UserCreate(ctx, userapi.AppNewUser{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*eerrs.Error)
				expResp := exp.(*eerrs.Error)

				return dbtest.CmpErrors(gotResp, expResp)
			},
		},
		{
			Name:    "badsig",
			Token:   sd.Admins[0].Token + "A",
			ExpResp: errs.Newf(http.StatusUnauthorized, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.UserCreate(ctx, userapi.AppNewUser{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*eerrs.Error)
				expResp := exp.(*eerrs.Error)

				return dbtest.CmpErrors(gotResp, expResp)
			},
		},
		{
			Name:    "wronguser",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_only]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := userapi.AppNewUser{
					Name:            "Bill Kennedy",
					Email:           "bill2@ardanlabs.com",
					Roles:           []string{"USER"},
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
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*eerrs.Error)
				expResp := exp.(*eerrs.Error)

				return dbtest.CmpErrors(gotResp, expResp)
			},
		},
	}

	return table
}
