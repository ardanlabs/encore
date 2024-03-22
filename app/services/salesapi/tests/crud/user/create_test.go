package user_test

import (
	"net/http"

	"github.com/ardanlabs/encore/app/services/salesapi/web/handlers/crud/usergrp"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func userCreate200(sd seedData) []tableData {
	table := []tableData{
		{
			name: "basic",
			//url:        "/v1/users",
			token: sd.admins[0].token,
			//statusCode: http.StatusCreated,
			// model: &usergrp.AppNewUser{
			// 	Name:            "Bill Kennedy",
			// 	Email:           "bill@ardanlabs.com",
			// 	Roles:           []string{"ADMIN"},
			// 	Department:      "IT",
			// 	Password:        "123",
			// 	PasswordConfirm: "123",
			// },
			expResp: &usergrp.AppUser{
				Name:       "Bill Kennedy",
				Email:      "bill@ardanlabs.com",
				Roles:      []string{"ADMIN"},
				Department: "IT",
				Enabled:    true,
			},
			cmpFunc: func(x interface{}, y interface{}) string {
				resp := x.(*usergrp.AppUser)
				expResp := y.(*usergrp.AppUser)

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

func userCreate400(sd seedData) []tableData {
	table := []tableData{
		{
			name: "missing-input",
			//url:        "/v1/users",
			token: sd.admins[0].token,
			//method:     http.MethodPost,
			//statusCode: http.StatusBadRequest,
			//model: &usergrp.AppNewUser{},
			expResp: toPointer(errs.NewResponse(http.StatusBadRequest, validate.FieldErrors{
				validate.FieldError{Field: "email", Err: "email is a required field"},
				validate.FieldError{Field: "name", Err: "name is a required field"},
				validate.FieldError{Field: "password", Err: "password is a required field"},
				validate.FieldError{Field: "roles", Err: "roles is a required field"},
			})),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			name: "bad-role",
			//url:        "/v1/users",
			token: sd.admins[0].token,
			//method:     http.MethodPost,
			//statusCode: http.StatusBadRequest,
			// model: &usergrp.AppNewUser{
			// 	Name:            "Bill Kennedy",
			// 	Email:           "bill@ardanlabs.com",
			// 	Roles:           []string{"BAD ROLE"},
			// 	Department:      "IT",
			// 	Password:        "123",
			// 	PasswordConfirm: "123",
			// },
			expResp: toPointer(errs.NewResponsef(http.StatusBadRequest, `parse: invalid role \"BAD ROLE\"`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func userCreate401(sd seedData) []tableData {
	table := []tableData{
		{
			name: "emptytoken",
			//url:        "/v1/users",
			token: "",
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			expResp: toPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			name: "badtoken",
			//url:        "/v1/users",
			token: sd.admins[0].token[:10],
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			expResp: toPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			name: "badsig",
			//url:        "/v1/users",
			token: sd.admins[0].token + "A",
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			expResp: toPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			name: "wronguser",
			//url:        "/v1/users",
			token: sd.users[0].token,
			//method:     http.MethodPost,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			expResp: toPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			cmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}
