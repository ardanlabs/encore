package home_test

import (
	"net/http"

	"github.com/ardanlabs/encore/app/services/salesapi/apis/crud/homeapi"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func homeCreate200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "basic",
			//url:        "/v1/homes",
			Token: sd.Users[0].Token,
			//method:     http.MethodPost,
			//statusCode: http.StatusCreated,
			// model: &homeapi.AppNewHome{
			// 	Type: "SINGLE FAMILY",
			// 	Address: homeapi.AppNewAddress{
			// 		Address1: "123 Mocking Bird Lane",
			// 		ZipCode:  "35810",
			// 		City:     "Huntsville",
			// 		State:    "AL",
			// 		Country:  "US",
			// 	},
			// },
			//// resp: &homeapi.AppHome{},
			ExpResp: &homeapi.AppHome{
				UserID: sd.Users[0].ID.String(),
				Type:   "SINGLE FAMILY",
				Address: homeapi.AppAddress{
					Address1: "123 Mocking Bird Lane",
					ZipCode:  "35810",
					City:     "Huntsville",
					State:    "AL",
					Country:  "US",
				},
			},
			CmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*homeapi.AppHome)
				expResp := y.(*homeapi.AppHome)

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

func homeCreate400(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "missing-input",
			//url:        "/v1/homes",
			Token: sd.Users[0].Token,
			//method:     http.MethodPost,
			//statusCode: http.StatusBadRequest,
			//model: &homeapi.AppNewHome{},
			//// resp:  &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponse(http.StatusBadRequest, validate.FieldErrors{
				validate.FieldError{Field: "address1", Err: "address1 is a required field"},
				validate.FieldError{Field: "city", Err: "city is a required field"},
				validate.FieldError{Field: "country", Err: "country is a required field"},
				validate.FieldError{Field: "state", Err: "state is a required field"},
				validate.FieldError{Field: "type", Err: "type is a required field"},
				validate.FieldError{Field: "zipCode", Err: "zipCode is a required field"},
			})),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "bad-type",
			//url:        "/v1/homes",
			Token: sd.Users[0].Token,
			//method:     http.MethodPost,
			//statusCode: http.StatusBadRequest,
			// model: &homeapi.AppNewHome{
			// 	Type: "BAD TYPE",
			// 	Address: homeapi.AppNewAddress{
			// 		Address1: "123 Mocking Bird Lane",
			// 		ZipCode:  "35810",
			// 		City:     "Huntsville",
			// 		State:    "AL",
			// 		Country:  "US",
			// 	},
			// },
			//// resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusBadRequest, `parse: invalid type \"BAD TYPE\"`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func homeCreate401(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "emptytoken",
			//url:        "/v1/homes",
			Token: "",
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			// resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "badtoken",
			//url:        "/v1/homes",
			Token: sd.Admins[0].Token[:10],
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			// resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "badsig",
			//url:        "/v1/homes",
			Token: sd.Admins[0].Token + "A",
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			// resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "wronguser",
			//url:        "/v1/homes",
			Token: sd.Admins[0].Token,
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			// resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
