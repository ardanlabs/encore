package user_test

import (
	"context"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/apis/services/sales"
	"github.com/ardanlabs/encore/app/api/apptest"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/app/core/crud/userapp"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func userCreateOk(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:  "basic",
			Token: sd.Admins[0].Token,
			ExpResp: userapp.User{
				Name:       "Bill Kennedy",
				Email:      "bill@ardanlabs.com",
				Roles:      []string{"ADMIN"},
				Department: "IT",
				Enabled:    true,
			},
			ExcFunc: func(ctx context.Context) any {
				app := userapp.NewUser{
					Name:            "Bill Kennedy",
					Email:           "bill@ardanlabs.com",
					Roles:           []string{"ADMIN"},
					Department:      "IT",
					Password:        "123",
					PasswordConfirm: "123",
				}

				resp, err := sales.UserCreate(ctx, app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(userapp.User)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(userapp.User)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func userCreateBad(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "missing",
			Token:   sd.Admins[0].Token,
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "validate: [{\"field\":\"name\",\"error\":\"name is a required field\"},{\"field\":\"email\",\"error\":\"email is a required field\"},{\"field\":\"roles\",\"error\":\"roles is a required field\"},{\"field\":\"password\",\"error\":\"password is a required field\"}]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.UserCreate(ctx, userapp.NewUser{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apptest.CmpAppErrors,
		},
		{
			Name:    "role",
			Token:   sd.Admins[0].Token,
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "parse: invalid role \"BAD ROLE\""),
			ExcFunc: func(ctx context.Context) any {
				app := userapp.NewUser{
					Name:            "Bill Kennedy",
					Email:           "bill2@ardanlabs.com",
					Roles:           []string{"BAD ROLE"},
					Department:      "IT",
					Password:        "123",
					PasswordConfirm: "123",
				}

				resp, err := sales.UserCreate(ctx, app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apptest.CmpAppErrors,
		},
	}

	return table
}

func userCreateAuth(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "&nbsp;",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.UserCreate(ctx, userapp.NewUser{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apptest.CmpAppErrors,
		},
		{
			Name:    "token",
			Token:   sd.Admins[0].Token[:10],
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.UserCreate(ctx, userapp.NewUser{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apptest.CmpAppErrors,
		},
		{
			Name:    "sig",
			Token:   sd.Admins[0].Token + "A",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.UserCreate(ctx, userapp.NewUser{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apptest.CmpAppErrors,
		},
		{
			Name:    "wronguser",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_only]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := userapp.NewUser{
					Name:            "Bill Kennedy",
					Email:           "bill2@ardanlabs.com",
					Roles:           []string{"USER"},
					Department:      "IT",
					Password:        "123",
					PasswordConfirm: "123",
				}

				resp, err := sales.UserCreate(ctx, app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apptest.CmpAppErrors,
		},
	}

	return table
}
