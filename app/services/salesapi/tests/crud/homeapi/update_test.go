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

func homeUpdate200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "basic",
			//url:        fmt.Sprintf("/v1/homes/%s", sd.Users[0].homes[0].ID),
			Token: sd.Users[0].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusOK,
			// model: &homeapi.AppUpdateHome{
			// 	Type: dbtest.StringPointer("SINGLE FAMILY"),
			// 	Address: &homeapi.AppUpdateAddress{
			// 		Address1: dbtest.StringPointer("123 Mocking Bird Lane"),
			// 		Address2: dbtest.StringPointer("apt 105"),
			// 		ZipCode:  dbtest.StringPointer("35810"),
			// 		City:     dbtest.StringPointer("Huntsville"),
			// 		State:    dbtest.StringPointer("AL"),
			// 		Country:  dbtest.StringPointer("US"),
			// 	},
			// },
			//resp: &homeapi.AppHome{},
			ExpResp: &homeapi.AppHome{
				UserID: sd.Users[0].ID.String(),
				Type:   "SINGLE FAMILY",
				Address: homeapi.AppAddress{
					Address1: "123 Mocking Bird Lane",
					Address2: "apt 105",
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

func homeUpdate400(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "bad-input",
			//url:        fmt.Sprintf("/v1/homes/%s", sd.Users[0].homes[0].ID),
			Token: sd.Users[0].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusBadRequest,
			// model: &homeapi.AppUpdateHome{
			// 	Address: &homeapi.AppUpdateAddress{
			// 		Address1: dbtest.StringPointer(""),
			// 		Address2: dbtest.StringPointer(""),
			// 		ZipCode:  dbtest.StringPointer(""),
			// 		City:     dbtest.StringPointer(""),
			// 		State:    dbtest.StringPointer(""),
			// 		Country:  dbtest.StringPointer(""),
			// 	},
			// },
			//resp: &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponse(http.StatusBadRequest, validate.FieldErrors{
				validate.FieldError{Field: "address1", Err: "address1 must be at least 1 character in length"},
				validate.FieldError{Field: "country", Err: "Key: 'AppUpdateHome.address.country' Error:Field validation for 'country' failed on the 'iso3166_1_alpha2' tag"},
				validate.FieldError{Field: "state", Err: "state must be at least 1 character in length"},
				validate.FieldError{Field: "zipCode", Err: "zipCode must be a valid numeric value"},
			})),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "bad-type",
			//url:        fmt.Sprintf("/v1/homes/%s", sd.Users[0].homes[0].ID),
			Token: sd.Users[0].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusBadRequest,
			// model: &homeapi.AppUpdateHome{
			// 	Type:    dbtest.StringPointer("BAD TYPE"),
			// 	Address: &homeapi.AppUpdateAddress{},
			// },
			//resp:    &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusBadRequest, `parse: invalid type \"BAD TYPE\"`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
	}

	return table
}

func homeUpdate401(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name: "emptytoken",
			//url:        fmt.Sprintf("/v1/homes/%s", sd.Users[0].homes[0].ID),
			Token: "",
			//method:     http.MethodPut,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "badsig",
			//url:        fmt.Sprintf("/v1/homes/%s", sd.Users[0].homes[0].ID),
			Token: sd.Users[0].Token + "A",
			//method:     http.MethodPut,
			//statusCode: http.StatusUnauthorized,
			//resp:       &middleware.Response{},
			ExpResp: dbtest.ToPointer(errs.NewResponsef(http.StatusUnauthorized, `Unauthorized`)),
			CmpFunc: func(x interface{}, y interface{}) string {
				return cmp.Diff(x, y)
			},
		},
		{
			Name: "wronguser",
			//url:        fmt.Sprintf("/v1/homes/%s", sd.admins[0].homes[0].ID),
			Token: sd.Users[0].Token,
			//method:     http.MethodPut,
			//statusCode: http.StatusUnauthorized,
			// model: &homeapi.AppUpdateHome{
			// 	Type: dbtest.StringPointer("SINGLE FAMILY"),
			// 	Address: &homeapi.AppUpdateAddress{
			// 		Address1: dbtest.StringPointer("123 Mocking Bird Lane"),
			// 		Address2: dbtest.StringPointer("apt 105"),
			// 		ZipCode:  dbtest.StringPointer("35810"),
			// 		City:     dbtest.StringPointer("Huntsville"),
			// 		State:    dbtest.StringPointer("AL"),
			// 		Country:  dbtest.StringPointer("US"),
			// 	},
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
