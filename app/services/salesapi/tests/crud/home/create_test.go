package home_test

import (
	"net/http"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/crud/homegrp"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func homeCreate200(sd seedData) []tableData {
	table := []tableData{
		{
			name:       "basic",
			url:        "/v1/homes",
			token:      sd.users[0].token,
			method:     http.MethodPost,
			statusCode: http.StatusCreated,
			model: &homegrp.AppNewHome{
				Type: "SINGLE FAMILY",
				Address: homegrp.AppNewAddress{
					Address1: "123 Mocking Bird Lane",
					ZipCode:  "35810",
					City:     "Huntsville",
					State:    "AL",
					Country:  "US",
				},
			},
			resp: &homegrp.AppHome{},
			expResp: &homegrp.AppHome{
				UserID: sd.users[0].ID.String(),
				Type:   "SINGLE FAMILY",
				Address: homegrp.AppAddress{
					Address1: "123 Mocking Bird Lane",
					ZipCode:  "35810",
					City:     "Huntsville",
					State:    "AL",
					Country:  "US",
				},
			},
			cmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*homegrp.AppHome)
				expResp := y.(*homegrp.AppHome)

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

func homeCreate400(sd seedData) []tableData {
	table := []tableData{
		{
			name:       "missing-input",
			url:        "/v1/homes",
			token:      sd.users[0].token,
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			model:      &homegrp.AppNewHome{},
			resp:       &middleware.Response{},
			expResp: toPointer(errs.NewResponse(http.StatusBadRequest, validate.FieldErrors{
				validate.FieldError{Field: "address1", Err: "address1 is a required field"},
				validate.FieldError{Field: "city", Err: "city is a required field"},
				validate.FieldError{Field: "country", Err: "country is a required field"},
				validate.FieldError{Field: "state", Err: "state is a required field"},
				validate.FieldError{Field: "type", Err: "type is a required field"},
				validate.FieldError{Field: "zipCode", Err: "zipCode is a required field"},
			})),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			name:       "bad-type",
			url:        "/v1/homes",
			token:      sd.users[0].token,
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			model: &homegrp.AppNewHome{
				Type: "BAD TYPE",
				Address: homegrp.AppNewAddress{
					Address1: "123 Mocking Bird Lane",
					ZipCode:  "35810",
					City:     "Huntsville",
					State:    "AL",
					Country:  "US",
				},
			},
			resp:    &middleware.Response{},
			expResp: toPointer(errs.NewResponsef(http.StatusBadRequest, `parse: invalid type \"BAD TYPE\"`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func homeCreate401(sd seedData) []tableData {
	table := []tableData{
		{
			name:       "emptytoken",
			url:        "/v1/homes",
			token:      "",
			method:     http.MethodPost,
			statusCode: http.StatusUnauthorized,
			resp:       &middleware.Response{},
			expResp:    toPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			name:       "badtoken",
			url:        "/v1/homes",
			token:      sd.admins[0].token[:10],
			method:     http.MethodPost,
			statusCode: http.StatusUnauthorized,
			resp:       &middleware.Response{},
			expResp:    toPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			name:       "badsig",
			url:        "/v1/homes",
			token:      sd.admins[0].token + "A",
			method:     http.MethodPost,
			statusCode: http.StatusUnauthorized,
			resp:       &middleware.Response{},
			expResp:    toPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			name:       "wronguser",
			url:        "/v1/homes",
			token:      sd.admins[0].token,
			method:     http.MethodPost,
			statusCode: http.StatusUnauthorized,
			resp:       &middleware.Response{},
			expResp:    toPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
