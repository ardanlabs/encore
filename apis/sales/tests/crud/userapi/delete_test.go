package user_test

import (
	"context"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/apis/sales"
	"github.com/ardanlabs/encore/app/api/apptest"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func userDeleteOk(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "user",
			Token:   sd.Users[1].Token,
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := sales.UserDelete(ctx, sd.Users[1].ID.String()); err != nil {
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
			Token:   sd.Admins[1].Token,
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := sales.UserDelete(ctx, sd.Admins[1].ID.String()); err != nil {
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

func userDeleteAuth(sd dbtest.SeedData) []apptest.AppTable {
	table := []apptest.AppTable{
		{
			Name:    "emptytoken",
			Token:   "",
			ExpResp: errs.Newf(eerrs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			ExcFunc: func(ctx context.Context) any {
				err := sales.UserDelete(ctx, "")
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
				err := sales.UserDelete(ctx, "")
				if err != nil {
					return err
				}

				return nil
			},
			CmpFunc: apptest.CmpAppErrors,
		},
		{
			Name:    "wronguser",
			Token:   sd.Users[1].Token,
			ExpResp: errs.Newf(eerrs.Unauthenticated, "user not enabled : query user: query: userID["+sd.Users[1].ID.String()+"]: db: user not found"),
			ExcFunc: func(ctx context.Context) any {
				err := sales.UserDelete(ctx, sd.Users[0].ID.String())
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
