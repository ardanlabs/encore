package product_test

import (
	"context"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/apis/services/salesapiweb"
	"github.com/ardanlabs/encore/app/api/apptest"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func productDeleteOk(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "user",
			Token:   sd.Users[0].Token,
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := salesapiweb.ProductDelete(ctx, sd.Users[0].Products[1].ID.String()); err != nil {
					return err
				}

				return nil
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "admin",
			Token:   sd.Admins[0].Token,
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := salesapiweb.ProductDelete(ctx, sd.Admins[0].Products[1].ID.String()); err != nil {
					return err
				}

				return nil
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func productDeleteAuth(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				err := salesapiweb.ProductDelete(ctx, "")
				if err != nil {
					return err
				}

				return nil
			},
			CmpFunc: apptest.CmpAppErrors,
		},
		{
			Name:    "sig",
			Token:   sd.Users[0].Token + "A",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				err := salesapiweb.ProductDelete(ctx, "")
				if err != nil {
					return err
				}

				return nil
			},
			CmpFunc: apptest.CmpAppErrors,
		},
		{
			Name:    "wronguser",
			Token:   sd.Users[0].Token,
			ExpResp: errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			ExcFunc: func(ctx context.Context) any {
				err := salesapiweb.ProductDelete(ctx, sd.Admins[0].Products[0].ID.String())
				if err != nil {
					return err
				}

				return nil
			},
			CmpFunc: apptest.CmpAppErrors,
		},
	}

	return table
}