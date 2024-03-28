package home_test

import (
	"context"
	"net/http"
	"time"

	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/homeapi"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func homeUpdateOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Users[0].Token,
			ExpResp: homeapi.AppHome{
				ID:     sd.Users[0].Homes[0].ID.String(),
				UserID: sd.Users[0].ID.String(),
				Type:   "SINGLE FAMILY",
				Address: homeapi.AppAddress{
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
				app := homeapi.AppUpdateHome{
					Type: dbtest.StringPointer("SINGLE FAMILY"),
					Address: &homeapi.AppUpdateAddress{
						Address1: dbtest.StringPointer("123 Mocking Bird Lane"),
						Address2: dbtest.StringPointer("apt 105"),
						ZipCode:  dbtest.StringPointer("35810"),
						City:     dbtest.StringPointer("Huntsville"),
						State:    dbtest.StringPointer("AL"),
						Country:  dbtest.StringPointer("US"),
					},
				}

				resp, err := salesapi.HomeUpdate(ctx, sd.Users[0].Homes[0].ID.String(), app)
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

func homeUpdateBad(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "input",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(http.StatusBadRequest, "validate: [{\"field\":\"address1\",\"error\":\"address1 must be at least 1 character in length\"},{\"field\":\"zipCode\",\"error\":\"zipCode must be a valid numeric value\"},{\"field\":\"state\",\"error\":\"state must be at least 1 character in length\"},{\"field\":\"country\",\"error\":\"Key: 'AppUpdateHome.address.country' Error:Field validation for 'country' failed on the 'iso3166_1_alpha2' tag\"}]"),
			ExcFunc: func(ctx context.Context) any {
				app := homeapi.AppUpdateHome{
					Address: &homeapi.AppUpdateAddress{
						Address1: dbtest.StringPointer(""),
						Address2: dbtest.StringPointer(""),
						ZipCode:  dbtest.StringPointer(""),
						City:     dbtest.StringPointer(""),
						State:    dbtest.StringPointer(""),
						Country:  dbtest.StringPointer(""),
					},
				}

				resp, err := salesapi.HomeUpdate(ctx, sd.Users[0].Homes[0].ID.String(), app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpAppErrors,
		},
		{
			Name:    "type",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(http.StatusBadRequest, "parse: invalid type \"BAD TYPE\""),
			ExcFunc: func(ctx context.Context) any {
				app := homeapi.AppUpdateHome{
					Type: dbtest.StringPointer("BAD TYPE"),
				}

				resp, err := salesapi.HomeUpdate(ctx, sd.Users[0].Homes[0].ID.String(), app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpAppErrors,
		},
	}

	return table
}

func homeUpdateAuth(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(http.StatusUnauthorized, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.HomeUpdate(ctx, "", homeapi.AppUpdateHome{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpAppErrors,
		},
		{
			Name:    "token",
			Token:   sd.Admins[0].Token[:10],
			ExpResp: errs.Newf(http.StatusUnauthorized, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.HomeUpdate(ctx, "", homeapi.AppUpdateHome{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpAppErrors,
		},
		{
			Name:    "sig",
			Token:   sd.Admins[0].Token + "A",
			ExpResp: errs.Newf(http.StatusUnauthorized, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.HomeUpdate(ctx, "", homeapi.AppUpdateHome{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpAppErrors,
		},
		{
			Name:    "wronguser",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := homeapi.AppUpdateHome{
					Type: dbtest.StringPointer("SINGLE FAMILY"),
					Address: &homeapi.AppUpdateAddress{
						Address1: dbtest.StringPointer("123 Mocking Bird Lane"),
						Address2: dbtest.StringPointer("apt 105"),
						ZipCode:  dbtest.StringPointer("35810"),
						City:     dbtest.StringPointer("Huntsville"),
						State:    dbtest.StringPointer("AL"),
						Country:  dbtest.StringPointer("US"),
					},
				}

				resp, err := salesapi.HomeUpdate(ctx, sd.Admins[0].Homes[0].ID.String(), app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpAppErrors,
		},
	}

	return table
}
