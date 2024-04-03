package user_test

import (
	"context"
	"time"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/apis/services/salesapiweb"
	"github.com/ardanlabs/encore/app/api/apptest"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/app/core/crud/userapp"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func userUpdateOk(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:  "basic",
			Token: sd.Users[0].Token,
			ExpResp: userapp.User{
				ID:          sd.Users[0].ID.String(),
				Name:        "Jack Kennedy",
				Email:       "jack@ardanlabs.com",
				Roles:       []string{"USER"},
				Department:  "IT",
				Enabled:     true,
				DateCreated: sd.Users[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Users[0].DateCreated.Format(time.RFC3339),
			},
			ExcFunc: func(ctx context.Context) any {
				app := userapp.UpdateUser{
					Name:            dbtest.StringPointer("Jack Kennedy"),
					Email:           dbtest.StringPointer("jack@ardanlabs.com"),
					Department:      dbtest.StringPointer("IT"),
					Password:        dbtest.StringPointer("123"),
					PasswordConfirm: dbtest.StringPointer("123"),
				}

				resp, err := salesapiweb.UserUpdate(ctx, sd.Users[0].ID.String(), app)
				if err != nil {
					return err
				}

				resp.DateUpdated = resp.DateCreated

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func userUpdateBad(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "input",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "validate: [{\"field\":\"email\",\"error\":\"email must be a valid email address\"},{\"field\":\"passwordConfirm\",\"error\":\"passwordConfirm must be equal to Password\"}]"),
			ExcFunc: func(ctx context.Context) any {
				app := userapp.UpdateUser{
					Email:           dbtest.StringPointer("jack@"),
					PasswordConfirm: dbtest.StringPointer("123"),
				}

				resp, err := salesapiweb.UserUpdate(ctx, sd.Users[0].ID.String(), app)
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
				app := userapp.UpdateUserRole{
					Roles: []string{"BAD ROLE"},
				}

				resp, err := salesapiweb.UserUpdateRole(ctx, sd.Admins[0].ID.String(), app)
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

func userUpdateAuth(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapiweb.UserUpdate(ctx, "", userapp.UpdateUser{})
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
				resp, err := salesapiweb.UserUpdate(ctx, sd.Admins[0].ID.String(), userapp.UpdateUser{})
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
				resp, err := salesapiweb.UserUpdate(ctx, sd.Admins[0].ID.String(), userapp.UpdateUser{})
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
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := userapp.UpdateUser{
					Name:            dbtest.StringPointer("Jack Kennedy"),
					Email:           dbtest.StringPointer("jack2@ardanlabs.com"),
					Department:      dbtest.StringPointer("IT"),
					Password:        dbtest.StringPointer("123"),
					PasswordConfirm: dbtest.StringPointer("123"),
				}

				resp, err := salesapiweb.UserUpdate(ctx, sd.Users[1].ID.String(), app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apptest.CmpAppErrors,
		},
		{
			Name:    "roleadminonly",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_only]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := userapp.UpdateUserRole{
					Roles: []string{"ADMIN"},
				}

				resp, err := salesapiweb.UserUpdateRole(ctx, sd.Users[1].ID.String(), app)
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