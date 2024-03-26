package product_test

import (
	"context"
	"net/http"

	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/data/dbtest"
)

func productDeleteOk(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "user",
			Token:   sd.Users[0].Token,
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := salesapi.ProductDelete(ctx, sd.Users[0].Products[1].ID.String()); err != nil {
					return err
				}

				return nil
			},
			CmpFunc: func(got any, exp any) string {
				if got != nil {
					return "error occurred"
				}

				return ""
			},
		},
		{
			Name:    "admin",
			Token:   sd.Admins[0].Token,
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := salesapi.ProductDelete(ctx, sd.Admins[0].Products[1].ID.String()); err != nil {
					return err
				}

				return nil
			},
			CmpFunc: func(got any, exp any) string {
				if got != nil {
					return "error occurred"
				}

				return ""
			},
		},
	}

	return table
}

func productDeleteAuth(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(http.StatusUnauthorized, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				err := salesapi.ProductDelete(ctx, "")
				if err != nil {
					return err
				}

				return nil
			},
			CmpFunc: dbtest.CmpErrors,
		},
		{
			Name:    "sig",
			Token:   sd.Users[0].Token + "A",
			ExpResp: errs.Newf(http.StatusUnauthorized, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				err := salesapi.ProductDelete(ctx, "")
				if err != nil {
					return err
				}

				return nil
			},
			CmpFunc: dbtest.CmpErrors,
		},
		{
			Name:    "wronguser",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				err := salesapi.ProductDelete(ctx, sd.Admins[0].Products[0].ID.String())
				if err != nil {
					return err
				}

				return nil
			},
			CmpFunc: dbtest.CmpErrors,
		},
	}

	return table
}
