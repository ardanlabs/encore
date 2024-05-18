package product_test

import (
	"context"
	"time"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/api/services/sales"
	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/domain/productapp"
	"github.com/ardanlabs/encore/app/sdk/errs"
	"github.com/ardanlabs/encore/business/sdk/dbtest"
	"github.com/google/go-cmp/cmp"
)

func updateOk(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:  "basic",
			Token: sd.Users[0].Token,
			ExpResp: productapp.Product{
				ID:          sd.Users[0].Products[0].ID.String(),
				UserID:      sd.Users[0].ID.String(),
				Name:        "Guitar",
				Cost:        10.34,
				Quantity:    10,
				DateCreated: sd.Users[0].Products[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Users[0].Products[0].DateCreated.Format(time.RFC3339),
			},
			ExcFunc: func(ctx context.Context) any {
				app := productapp.UpdateProduct{
					Name:     dbtest.StringPointer("Guitar"),
					Cost:     dbtest.FloatPointer(10.34),
					Quantity: dbtest.IntPointer(10),
				}

				resp, err := sales.ProductUpdate(ctx, sd.Users[0].Products[0].ID.String(), app)
				if err != nil {
					return err
				}

				resp.DateUpdated = resp.DateCreated

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				if _, exists := got.(productapp.Product); !exists {
					return "error occurred"
				}

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
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "validate: [{\"field\":\"cost\",\"error\":\"cost must be 0 or greater\"},{\"field\":\"quantity\",\"error\":\"quantity must be 1 or greater\"}]"),
			ExcFunc: func(ctx context.Context) any {
				app := productapp.UpdateProduct{
					Cost:     dbtest.FloatPointer(-10.34),
					Quantity: dbtest.IntPointer(-10),
				}

				resp, err := sales.ProductUpdate(ctx, sd.Users[0].ID.String(), app)
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
				resp, err := sales.ProductUpdate(ctx, "", productapp.UpdateProduct{})
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
				resp, err := sales.ProductUpdate(ctx, "", productapp.UpdateProduct{})
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
				resp, err := sales.ProductUpdate(ctx, "", productapp.UpdateProduct{})
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
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := productapp.UpdateProduct{
					Name:     dbtest.StringPointer("Guitar"),
					Cost:     dbtest.FloatPointer(10.34),
					Quantity: dbtest.IntPointer(10),
				}

				resp, err := sales.ProductUpdate(ctx, sd.Admins[0].Products[0].ID.String(), app)
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
