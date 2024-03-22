package product_test

import (
	"fmt"
	"net/http"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/crud/productgrp"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func productUpdate200(sd seedData) []tableData {
	table := []tableData{
		{
			name:       "basic",
			url:        fmt.Sprintf("/v1/products/%s", sd.users[1].products[0].ID),
			token:      sd.users[1].token,
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			model: &productgrp.AppUpdateProduct{
				Name:     dbtest.StringPointer("Guitar"),
				Cost:     dbtest.FloatPointer(10.34),
				Quantity: dbtest.IntPointer(10),
			},
			resp: &productgrp.AppProduct{},
			expResp: &productgrp.AppProduct{
				Name:     "Guitar",
				UserID:   sd.users[1].ID.String(),
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

func productUpdate400(sd seedData) []tableData {
	table := []tableData{
		{
			name:       "bad-input",
			url:        fmt.Sprintf("/v1/products/%s", sd.users[1].products[0].ID),
			token:      sd.users[1].token,
			method:     http.MethodPut,
			statusCode: http.StatusBadRequest,
			model: &productgrp.AppUpdateProduct{
				Cost:     dbtest.FloatPointer(-1.0),
				Quantity: dbtest.IntPointer(0),
			},
			resp: &middleware.Response{},
			expResp: toPointer(errs.NewResponse(http.StatusBadRequest, validate.FieldErrors{
				validate.FieldError{Field: "cost", Err: "cost must be 0 or greate"},
				validate.FieldError{Field: "quantity", Err: "quantity must be 1 or greater"},
			})),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func productUpdate401(sd seedData) []tableData {
	table := []tableData{
		{
			name:       "emptytoken",
			url:        fmt.Sprintf("/v1/products/%s", sd.users[1].products[0].ID),
			token:      "",
			method:     http.MethodPut,
			statusCode: http.StatusUnauthorized,
			resp:       &middleware.Response{},
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			name:       "badsig",
			url:        fmt.Sprintf("/v1/products/%s", sd.users[1].products[0].ID),
			token:      sd.users[1].token + "A",
			method:     http.MethodPut,
			statusCode: http.StatusUnauthorized,
			expResp:    toPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			name:       "wronguser",
			url:        fmt.Sprintf("/v1/products/%s", sd.admins[1].products[0].ID),
			token:      sd.users[0].token,
			method:     http.MethodPut,
			statusCode: http.StatusUnauthorized,
			model: &productgrp.AppUpdateProduct{
				Name:     dbtest.StringPointer("Guitar"),
				Cost:     dbtest.FloatPointer(10.34),
				Quantity: dbtest.IntPointer(10),
			},
			resp:    &middleware.Response{},
			expResp: toPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
