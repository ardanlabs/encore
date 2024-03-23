package product_test

import (
	"net/http"

	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/crud/productgrp"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func productUpdate200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "basic",
			//url:        fmt.Sprintf("/v1/products/%s", sd.Users[1].products[0].ID),
			Token: sd.Users[1].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusOK,
			// model: &productgrp.AppUpdateProduct{
			// 	Name:     dbtest.StringPointer("Guitar"),
			// 	Cost:     dbtest.FloatPointer(10.34),
			// 	Quantity: dbtest.IntPointer(10),
			// },
			//resp: &productgrp.AppProduct{},
			ExpResp: &productgrp.AppProduct{
				Name:     "Guitar",
				UserID:   sd.Users[1].ID.String(),
				Cost:     10.34,
				Quantity: 10,
			},
			CmpFunc: func(x interface{}, y interface{}) string {
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

func productUpdate400(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "bad-input",
			//url:        fmt.Sprintf("/v1/products/%s", sd.Users[1].products[0].ID),
			Token: sd.Users[1].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusBadRequest,
			// model: &productgrp.AppUpdateProduct{
			// 	Cost:     dbtest.FloatPointer(-1.0),
			// 	Quantity: dbtest.IntPointer(0),
			// },
			//resp: &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponse(http.StatusBadRequest, validate.FieldErrors{
				validate.FieldError{Field: "cost", Err: "cost must be 0 or greate"},
				validate.FieldError{Field: "quantity", Err: "quantity must be 1 or greater"},
			})),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func productUpdate401(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "emptytoken",
			//url:        fmt.Sprintf("/v1/products/%s", sd.Users[1].products[0].ID),
			Token: "",
			//method:     http.MethodPut,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "badsig",
			//url:        fmt.Sprintf("/v1/products/%s", sd.Users[1].products[0].ID),
			Token: sd.Users[1].Token + "A",
			//method:     http.MethodPut,
			//statusCode: http.StatusUnauthorized,
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "wronguser",
			//url:        fmt.Sprintf("/v1/products/%s", sd.admins[1].products[0].ID),
			Token: sd.Users[0].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusUnauthorized,
			// model: &productgrp.AppUpdateProduct{
			// 	Name:     dbtest.StringPointer("Guitar"),
			// 	Cost:     dbtest.FloatPointer(10.34),
			// 	Quantity: dbtest.IntPointer(10),
			// },
			//resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
