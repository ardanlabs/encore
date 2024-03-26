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

func userUpdateOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Users[0].Token,
			ExpResp: userapi.AppUser{
				Name:       "Jack Kennedy",
				Email:      "jack@ardanlabs.com",
				Roles:      []string{"ADMIN"},
				Department: "IT",
				Enabled:    true,
			},
			ExcFunc: func(ctx context.Context) any {
				app := userapi.AppUpdateUser{
					Name:            dbtest.StringPointer("Jack Kennedy"),
					Email:           dbtest.StringPointer("jack@ardanlabs.com"),
					Roles:           []string{"ADMIN"},
					Department:      dbtest.StringPointer("IT"),
					Password:        dbtest.StringPointer("123"),
					PasswordConfirm: dbtest.StringPointer("123"),
				}

				resp, err := salesapi.UserUpdate(ctx, sd.Users[0].ID.String(), app)
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

func userUpdateBad(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "bad-input",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(http.StatusBadRequest, "validate: [{\"field\":\"email\",\"error\":\"email must be a valid email address\"},{\"field\":\"passwordConfirm\",\"error\":\"passwordConfirm must be equal to Password\"}]"),
			ExcFunc: func(ctx context.Context) any {
				app := userapi.AppUpdateUser{
					Email:           dbtest.StringPointer("jack@"),
					PasswordConfirm: dbtest.StringPointer("123"),
				}

				resp, err := salesapi.UserUpdate(ctx, sd.Users[0].ID.String(), app)
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
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(http.StatusBadRequest, "parse: invalid role \"BAD ROLE\""),
			ExcFunc: func(ctx context.Context) any {
				app := userapi.AppUpdateUser{
					Roles: []string{"BAD ROLE"},
				}

				resp, err := salesapi.UserUpdate(ctx, sd.Users[0].ID.String(), app)
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

func userUpdateAuth(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(http.StatusUnauthorized, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.UserUpdate(ctx, "", userapi.AppUpdateUser{})
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
				resp, err := salesapi.UserUpdate(ctx, sd.Admins[0].ID.String(), userapi.AppUpdateUser{})
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
				resp, err := salesapi.UserUpdate(ctx, sd.Admins[0].ID.String(), userapi.AppUpdateUser{})
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
			ExpResp: errs.Newf(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := userapi.AppUpdateUser{
					Name:            dbtest.StringPointer("Jack Kennedy"),
					Email:           dbtest.StringPointer("jack2@ardanlabs.com"),
					Roles:           []string{"ADMIN"},
					Department:      dbtest.StringPointer("IT"),
					Password:        dbtest.StringPointer("123"),
					PasswordConfirm: dbtest.StringPointer("123"),
				}

				resp, err := salesapi.UserUpdate(ctx, sd.Users[1].ID.String(), app)
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
