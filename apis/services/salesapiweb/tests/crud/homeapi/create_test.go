package home_test

import (
	"context"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/apis/services/salesapiweb"
	"github.com/ardanlabs/encore/app/api/apptest"
	"github.com/ardanlabs/encore/app/api/errs"
	homeapp "github.com/ardanlabs/encore/app/core/crud/homeapp"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func homeCreateOk(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:  "basic",
			Token: sd.Users[0].Token,
			ExpResp: homeapp.Home{
				UserID: sd.Users[0].ID.String(),
				Type:   "SINGLE FAMILY",
				Address: homeapp.Address{
					Address1: "123 Mocking Bird Lane",
					ZipCode:  "35810",
					City:     "Huntsville",
					State:    "AL",
					Country:  "US",
				},
			},
			ExcFunc: func(ctx context.Context) any {
				app := homeapp.NewHome{
					Type: "SINGLE FAMILY",
					Address: homeapp.NewAddress{
						Address1: "123 Mocking Bird Lane",
						ZipCode:  "35810",
						City:     "Huntsville",
						State:    "AL",
						Country:  "US",
					},
				}

				resp, err := salesapiweb.HomeCreate(ctx, app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(homeapp.Home)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(homeapp.Home)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func homeCreateBad(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "missing",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "validate: [{\"field\":\"type\",\"error\":\"type is a required field\"},{\"field\":\"address1\",\"error\":\"address1 is a required field\"},{\"field\":\"zipCode\",\"error\":\"zipCode is a required field\"},{\"field\":\"city\",\"error\":\"city is a required field\"},{\"field\":\"state\",\"error\":\"state is a required field\"},{\"field\":\"country\",\"error\":\"country is a required field\"}]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapiweb.HomeCreate(ctx, homeapp.NewHome{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apptest.CmpAppErrors,
		},
		{
			Name:    "type",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "parse: invalid type \"BAD TYPE\""),
			ExcFunc: func(ctx context.Context) any {
				app := homeapp.NewHome{
					Type: "BAD TYPE",
					Address: homeapp.NewAddress{
						Address1: "123 Mocking Bird Lane",
						ZipCode:  "35810",
						City:     "Huntsville",
						State:    "AL",
						Country:  "US",
					},
				}

				resp, err := salesapiweb.HomeCreate(ctx, app)
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

func homeCreateAuth(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapiweb.HomeCreate(ctx, homeapp.NewHome{})
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
				resp, err := salesapiweb.HomeCreate(ctx, homeapp.NewHome{})
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
				resp, err := salesapiweb.HomeCreate(ctx, homeapp.NewHome{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: apptest.CmpAppErrors,
		},
		{
			Name:    "wronguser",
			Token:   sd.Admins[0].Token,
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[[{ADMIN}]] rule[rule_user_only]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := homeapp.NewHome{
					Type: "SINGLE FAMILY",
					Address: homeapp.NewAddress{
						Address1: "123 Mocking Bird Lane",
						ZipCode:  "35810",
						City:     "Huntsville",
						State:    "AL",
						Country:  "US",
					},
				}

				resp, err := salesapiweb.HomeCreate(ctx, app)
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
