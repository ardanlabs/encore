package tran_test

import (
	"context"

	"github.com/ardanlabs/encore/api/services/sales"
	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/domain/tranapp"
	"github.com/google/go-cmp/cmp"
)

func createOk(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:  "basic",
			Token: sd.Admins[0].Token,
			ExpResp: tranapp.Product{
				Name:     "Guitar",
				Cost:     10.34,
				Quantity: 10,
			},
			ExcFunc: func(ctx context.Context) any {
				app := tranapp.NewTran{
					Product: tranapp.NewProduct{
						Name:     "Guitar",
						Cost:     10.34,
						Quantity: 10,
					},
					User: tranapp.NewUser{
						Name:            "Bill Kennedy",
						Email:           "bill@ardanlabs.com",
						Roles:           []string{"ADMIN"},
						Department:      "IT",
						Password:        "123",
						PasswordConfirm: "123",
					},
				}

				resp, err := sales.TranCreate(ctx, app)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(tranapp.Product)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(tranapp.Product)

				expResp.ID = gotResp.ID
				expResp.UserID = gotResp.UserID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
