package home_test

import (
	"context"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/services/salesapi"
	homeapp "github.com/ardanlabs/encore/app/services/salesapi/core/crud/homeapp"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func homeCreateOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Users[0].Token,
			ExpResp: homeapp.AppHome{
				UserID: sd.Users[0].ID.String(),
				Type:   "SINGLE FAMILY",
				Address: homeapp.AppAddress{
					Address1: "123 Mocking Bird Lane",
					ZipCode:  "35810",
					City:     "Huntsville",
					State:    "AL",
					Country:  "US",
				},
			},
			ExcFunc: func(ctx context.Context) any {
				app := homeapp.AppNewHome{
					Type: "SINGLE FAMILY",
					Address: homeapp.AppNewAddress{
						Address1: "123 Mocking Bird Lane",
						ZipCode:  "35810",
						City:     "Huntsville",
						State:    "AL",
						Country:  "US",
					},
				}

				resp, err := salesapi.HomeCreate(ctx, app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(homeapp.AppHome)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(homeapp.AppHome)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func homeCreateBad(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "missing",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "validate: [{\"field\":\"type\",\"error\":\"type is a required field\"},{\"field\":\"address1\",\"error\":\"address1 is a required field\"},{\"field\":\"zipCode\",\"error\":\"zipCode is a required field\"},{\"field\":\"city\",\"error\":\"city is a required field\"},{\"field\":\"state\",\"error\":\"state is a required field\"},{\"field\":\"country\",\"error\":\"country is a required field\"}]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.HomeCreate(ctx, homeapp.AppNewHome{})
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
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "parse: invalid type \"BAD TYPE\""),
			ExcFunc: func(ctx context.Context) any {
				app := homeapp.AppNewHome{
					Type: "BAD TYPE",
					Address: homeapp.AppNewAddress{
						Address1: "123 Mocking Bird Lane",
						ZipCode:  "35810",
						City:     "Huntsville",
						State:    "AL",
						Country:  "US",
					},
				}

				resp, err := salesapi.HomeCreate(ctx, app)
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

func homeCreateAuth(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.HomeCreate(ctx, homeapp.AppNewHome{})
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
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.HomeCreate(ctx, homeapp.AppNewHome{})
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
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.HomeCreate(ctx, homeapp.AppNewHome{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpAppErrors,
		},
		{
			Name:    "wronguser",
			Token:   sd.Admins[0].Token,
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[[{ADMIN}]] rule[rule_user_only]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := homeapp.AppNewHome{
					Type: "SINGLE FAMILY",
					Address: homeapp.AppNewAddress{
						Address1: "123 Mocking Bird Lane",
						ZipCode:  "35810",
						City:     "Huntsville",
						State:    "AL",
						Country:  "US",
					},
				}

				resp, err := salesapi.HomeCreate(ctx, app)
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
