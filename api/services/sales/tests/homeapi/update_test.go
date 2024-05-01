package home_test

import (
	"context"
	"time"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/api/services/sales"
	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/app/domain/homeapp"
	"github.com/ardanlabs/encore/business/api/dbtest"
	"github.com/google/go-cmp/cmp"
)

func updateOk(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:  "basic",
			Token: sd.Users[0].Token,
			ExpResp: homeapp.Home{
				ID:     sd.Users[0].Homes[0].ID.String(),
				UserID: sd.Users[0].ID.String(),
				Type:   "SINGLE FAMILY",
				Address: homeapp.Address{
					Address1: "123 Mocking Bird Lane",
					Address2: "apt 105",
					ZipCode:  "35810",
					City:     "Huntsville",
					State:    "AL",
					Country:  "US",
				},
				DateCreated: sd.Users[0].Homes[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Users[0].Homes[0].DateCreated.Format(time.RFC3339),
			},
			ExcFunc: func(ctx context.Context) any {
				app := homeapp.UpdateHome{
					Type: dbtest.StringPointer("SINGLE FAMILY"),
					Address: &homeapp.UpdateAddress{
						Address1: dbtest.StringPointer("123 Mocking Bird Lane"),
						Address2: dbtest.StringPointer("apt 105"),
						ZipCode:  dbtest.StringPointer("35810"),
						City:     dbtest.StringPointer("Huntsville"),
						State:    dbtest.StringPointer("AL"),
						Country:  dbtest.StringPointer("US"),
					},
				}

				resp, err := sales.HomeUpdate(ctx, sd.Users[0].Homes[0].ID.String(), app)
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

func updateBad(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:    "input",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "validate: [{\"field\":\"address1\",\"error\":\"address1 must be at least 1 character in length\"},{\"field\":\"zipCode\",\"error\":\"zipCode must be a valid numeric value\"},{\"field\":\"state\",\"error\":\"state must be at least 1 character in length\"},{\"field\":\"country\",\"error\":\"Key: 'UpdateHome.address.country' Error:Field validation for 'country' failed on the 'iso3166_1_alpha2' tag\"}]"),
			ExcFunc: func(ctx context.Context) any {
				app := homeapp.UpdateHome{
					Address: &homeapp.UpdateAddress{
						Address1: dbtest.StringPointer(""),
						Address2: dbtest.StringPointer(""),
						ZipCode:  dbtest.StringPointer(""),
						City:     dbtest.StringPointer(""),
						State:    dbtest.StringPointer(""),
						Country:  dbtest.StringPointer(""),
					},
				}

				resp, err := sales.HomeUpdate(ctx, sd.Users[0].Homes[0].ID.String(), app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apitest.CmpAppErrors,
		},
		{
			Name:    "type",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "parse: invalid type \"BAD TYPE\""),
			ExcFunc: func(ctx context.Context) any {
				app := homeapp.UpdateHome{
					Type: dbtest.StringPointer("BAD TYPE"),
				}

				resp, err := sales.HomeUpdate(ctx, sd.Users[0].Homes[0].ID.String(), app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apitest.CmpAppErrors,
		},
	}

	return table
}

func updateAuth(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:    "emptytoken",
			Token:   "&nbsp;",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.HomeUpdate(ctx, "", homeapp.UpdateHome{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apitest.CmpAppErrors,
		},
		{
			Name:    "token",
			Token:   sd.Admins[0].Token[:10],
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.HomeUpdate(ctx, "", homeapp.UpdateHome{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apitest.CmpAppErrors,
		},
		{
			Name:    "sig",
			Token:   sd.Admins[0].Token + "A",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.HomeUpdate(ctx, "", homeapp.UpdateHome{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apitest.CmpAppErrors,
		},
		{
			Name:    "wronguser",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := homeapp.UpdateHome{
					Type: dbtest.StringPointer("SINGLE FAMILY"),
					Address: &homeapp.UpdateAddress{
						Address1: dbtest.StringPointer("123 Mocking Bird Lane"),
						Address2: dbtest.StringPointer("apt 105"),
						ZipCode:  dbtest.StringPointer("35810"),
						City:     dbtest.StringPointer("Huntsville"),
						State:    dbtest.StringPointer("AL"),
						Country:  dbtest.StringPointer("US"),
					},
				}

				resp, err := sales.HomeUpdate(ctx, sd.Admins[0].Homes[0].ID.String(), app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apitest.CmpAppErrors,
		},
	}

	return table
}
