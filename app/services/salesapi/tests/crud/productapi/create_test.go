package product_test

import (
	"context"
	"net/http"

	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/productapi"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func productCreateOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Users[0].Token,
			ExpResp: productapi.AppProduct{
				UserID:   sd.Users[0].ID.String(),
				Name:     "Guitar",
				Cost:     10.34,
				Quantity: 10,
			},
			ExcFunc: func(ctx context.Context) any {
				app := productapi.AppNewProduct{
					Name:     "Guitar",
					Cost:     10.34,
					Quantity: 10,
				}

				resp, err := salesapi.ProductCreate(ctx, app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(productapi.AppProduct)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(productapi.AppProduct)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func productCreateBad(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "missing",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(http.StatusBadRequest, "validate: [{\"field\":\"name\",\"error\":\"name is a required field\"},{\"field\":\"cost\",\"error\":\"cost is a required field\"},{\"field\":\"quantity\",\"error\":\"quantity is a required field\"}]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.ProductCreate(ctx, productapi.AppNewProduct{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpErrors,
		},
	}

	return table
}

func productCreateAuth(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(http.StatusUnauthorized, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.ProductCreate(ctx, productapi.AppNewProduct{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpErrors,
		},
		{
			Name:    "token",
			Token:   sd.Admins[0].Token[:10],
			ExpResp: errs.Newf(http.StatusUnauthorized, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.ProductCreate(ctx, productapi.AppNewProduct{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpErrors,
		},
		{
			Name:    "sig",
			Token:   sd.Admins[0].Token + "A",
			ExpResp: errs.Newf(http.StatusUnauthorized, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.ProductCreate(ctx, productapi.AppNewProduct{})
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpErrors,
		},
		{
			Name:    "wronguser",
			Token:   sd.Admins[0].Token,
			ExpResp: errs.Newf(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[[{ADMIN}]] rule[rule_user_only]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := productapi.AppNewProduct{
					Name:     "Guitar",
					Cost:     10.34,
					Quantity: 10,
				}

				resp, err := salesapi.ProductCreate(ctx, app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: dbtest.CmpErrors,
		},
	}

	return table
}
