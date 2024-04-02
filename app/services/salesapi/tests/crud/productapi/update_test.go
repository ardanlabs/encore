package product_test

import (
	"context"
	"time"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/productapi"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func productUpdateOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:  "basic",
			Token: sd.Users[0].Token,
			ExpResp: productapi.AppProduct{
				ID:          sd.Users[0].Products[0].ID.String(),
				UserID:      sd.Users[0].ID.String(),
				Name:        "Guitar",
				Cost:        10.34,
				Quantity:    10,
				DateCreated: sd.Users[0].Products[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Users[0].Products[0].DateCreated.Format(time.RFC3339),
			},
			ExcFunc: func(ctx context.Context) any {
				app := productapi.AppUpdateProduct{
					Name:     dbtest.StringPointer("Guitar"),
					Cost:     dbtest.FloatPointer(10.34),
					Quantity: dbtest.IntPointer(10),
				}

				resp, err := salesapi.ProductUpdate(ctx, sd.Users[0].Products[0].ID.String(), app)
				if err != nil {
					return err
				}

				resp.DateUpdated = resp.DateCreated

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				if _, exists := got.(productapi.AppProduct); !exists {
					return "error occurred"
				}

				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func productUpdateBad(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "input",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.FailedPrecondition, "validate: [{\"field\":\"cost\",\"error\":\"cost must be 0 or greater\"},{\"field\":\"quantity\",\"error\":\"quantity must be 1 or greater\"}]"),
			ExcFunc: func(ctx context.Context) any {
				app := productapi.AppUpdateProduct{
					Cost:     dbtest.FloatPointer(-10.34),
					Quantity: dbtest.IntPointer(-10),
				}

				resp, err := salesapi.ProductUpdate(ctx, sd.Users[0].ID.String(), app)
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

func productUpdateAuth(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				resp, err := salesapi.ProductUpdate(ctx, "", productapi.AppUpdateProduct{})
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
				resp, err := salesapi.ProductUpdate(ctx, "", productapi.AppUpdateProduct{})
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
				resp, err := salesapi.ProductUpdate(ctx, "", productapi.AppUpdateProduct{})
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
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				app := productapi.AppUpdateProduct{
					Name:     dbtest.StringPointer("Guitar"),
					Cost:     dbtest.FloatPointer(10.34),
					Quantity: dbtest.IntPointer(10),
				}

				resp, err := salesapi.ProductUpdate(ctx, sd.Admins[0].Products[0].ID.String(), app)
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
