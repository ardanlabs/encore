package product_test

import (
	"context"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/apis/sales"
	"github.com/ardanlabs/encore/app/api/apptest"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/app/core/crud/productapp"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func productCreateOk(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:  "basic",
			Token: sd.Users[0].Token,
			ExpResp: productapp.Product{
				UserID:   sd.Users[0].ID.String(),
				Name:     "Guitar",
				Cost:     10.34,
				Quantity: 10,
			},
			ExcFunc: func(ctx context.Context) any {
				app := productapp.NewProduct{
					Name:     "Guitar",
					Cost:     10.34,
					Quantity: 10,
				}

				resp, err := sales.ProductCreate(ctx, app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(productapp.Product)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(productapp.Product)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func productCreateBad(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "missing",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "validate: [{\"field\":\"name\",\"error\":\"name is a required field\"},{\"field\":\"cost\",\"error\":\"cost is a required field\"},{\"field\":\"quantity\",\"error\":\"quantity is a required field\"}]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.ProductCreate(ctx, productapp.NewProduct{})
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

func productCreateAuth(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := sales.ProductCreate(ctx, productapp.NewProduct{})
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
				resp, err := sales.ProductCreate(ctx, productapp.NewProduct{})
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
				resp, err := sales.ProductCreate(ctx, productapp.NewProduct{})
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
				app := productapp.NewProduct{
					Name:     "Guitar",
					Cost:     10.34,
					Quantity: 10,
				}

				resp, err := sales.ProductCreate(ctx, app)
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
