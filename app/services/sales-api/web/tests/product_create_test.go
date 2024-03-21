package tests

import (
	"net/http"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/services/sales-api/web/handlers/productgrp"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func productCreate200(sd seedData) []tableData {
	table := []tableData{
		{
			name:       "basic",
			url:        "/v1/products",
			token:      sd.users[0].token,
			method:     http.MethodPost,
			statusCode: http.StatusCreated,
			model: &productgrp.AppNewProduct{
				Name:     "Guitar",
				Cost:     10.34,
				Quantity: 10,
			},
			resp: &productgrp.AppProduct{},
			expResp: &productgrp.AppProduct{
				Name:     "Guitar",
				UserID:   sd.users[0].ID.String(),
				Cost:     10.34,
				Quantity: 10,
			},
			cmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*productgrp.AppProduct)
				expResp := y.(*productgrp.AppProduct)

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

func productCreate400(sd seedData) []tableData {
	table := []tableData{
		{
			name:       "missing-input",
			url:        "/v1/products",
			token:      sd.users[0].token,
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			model:      &productgrp.AppNewProduct{},
			resp:       &middleware.Response{},
			expResp: toPointer(errs.NewResponse(http.StatusBadRequest, validate.FieldErrors{
				validate.FieldError{Field: "cost", Err: "cost is a required field"},
				validate.FieldError{Field: "name", Err: "name is a required field"},
				validate.FieldError{Field: "quantity", Err: "quantity is a required field"},
			})),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func productCreate401(sd seedData) []tableData {
	table := []tableData{
		{
			name:       "emptytoken",
			url:        "/v1/products",
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
			url:        "/v1/products",
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
			url:        "/v1/products",
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
			url:        "/v1/products",
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
