package product_test

import (
	"net/http"

	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/productapi"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func productCreate200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "basic",
			//url:        "/v1/products",
			Token: sd.Users[0].Token,
			//method:     http.MethodPost,
			//statusCode: http.StatusCreated,
			// model: &productapi.AppNewProduct{
			// 	Name:     "Guitar",
			// 	Cost:     10.34,
			// 	Quantity: 10,
			// },
			//resp: &productapi.AppProduct{},
			ExpResp: &productapi.AppProduct{
				Name:     "Guitar",
				UserID:   sd.Users[0].ID.String(),
				Cost:     10.34,
				Quantity: 10,
			},
			CmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*productapi.AppProduct)
				expResp := y.(*productapi.AppProduct)

				if _, err := uuid.Parse(resp.ID); err != nil {
					return "bad uuid for ID"
				}

				if resp.DateCreated == "" {
					return "missing date created"
				}

				if resp.DateUpdated == "" {
					return "missing date updated"
				}

				expResp.ID = resp.ID
				expResp.DateCreated = resp.DateCreated
				expResp.DateUpdated = resp.DateUpdated

				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func productCreate400(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "missing-input",
			//url:        "/v1/products",
			Token: sd.Users[0].Token,
			//method:     http.MethodPost,
			//statusCode: http.StatusBadRequest,
			//model: &productapi.AppNewProduct{},
			//resp:       &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponse(http.StatusBadRequest, validate.FieldErrors{
				validate.FieldError{Field: "cost", Err: "cost is a required field"},
				validate.FieldError{Field: "name", Err: "name is a required field"},
				validate.FieldError{Field: "quantity", Err: "quantity is a required field"},
			})),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func productCreate401(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "emptytoken",
			//url:        "/v1/products",
			Token: "",
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "badtoken",
			//url:        "/v1/products",
			Token: sd.Admins[0].Token[:10],
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "badsig",
			//url:        "/v1/products",
			Token: sd.Admins[0].Token + "A",
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "wronguser",
			//url:        "/v1/products",
			Token: sd.Admins[0].Token,
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
