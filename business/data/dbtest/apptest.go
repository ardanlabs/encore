package dbtest

import (
	"context"
	"testing"

	eauth "encore.dev/beta/auth"
	eerrs "encore.dev/beta/errs"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/api/mid"
	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

// AppTable represent fields needed for running an app test.
type AppTable struct {
	Name    string
	Token   string
	ExpResp any
	ExcFunc func(ctx context.Context) any
	CmpFunc func(got any, exp any) string
}

// User represents an app user specified for the test.
type User struct {
	user.User
	Token    string
	Products []product.Product
	Homes    []home.Home
}

// SeedData represents data that was seeded for the test.
type SeedData struct {
	Users  []User
	Admins []User
}

// ToPointer converts a middleware reponose value to a pointer.
func ToPointer(r middleware.Response) *middleware.Response {
	return &r
}

// Service defines the method set required to exist for any encore service type.
type Service interface {
	AuthHandler(ctx context.Context, ap *mid.AuthParams) (eauth.UID, *auth.Claims, error)
}

// =============================================================================

// AppTest contains functions for executing an app test.
type AppTest struct {
	Service Service
}

// Test performs the actual test logic based on the table data.
func (at *AppTest) Test(t *testing.T, table []AppTable, testName string) {
	log := func(got any, exp any) {
		t.Log("GOT")
		t.Logf("%#v", got)
		t.Log("EXP")
		t.Logf("%#v", exp)
		t.Fatalf("Should get the expected response")
	}

	for _, tt := range table {
		f := func(t *testing.T) {
			ctx := context.Background()

			t.Log("Calling authHandler")
			ctx, err := at.authHandler(ctx, tt.Token)
			if err != nil {
				diff := tt.CmpFunc(err, tt.ExpResp)
				if diff != "" {
					log(err, tt.ExpResp)
				}
				return
			}

			t.Log("Calling excFunc")
			got := tt.ExcFunc(ctx)

			diff := tt.CmpFunc(got, tt.ExpResp)
			if diff != "" {
				log(got, tt.ExpResp)
			}
		}

		t.Run(testName+"-"+tt.Name, f)
	}
}

func (at *AppTest) authHandler(ctx context.Context, token string) (context.Context, error) {
	uid, claims, err := at.Service.AuthHandler(ctx, &mid.AuthParams{
		Authorization: "Bearer " + token,
	})

	if err != nil {
		return ctx, err
	}

	return eauth.WithContext(ctx, uid, claims), nil
}

// =============================================================================

// CmpErrors compares two encore error values. If they are not equal, the
// reason is returned.
func CmpErrors(got any, exp any) string {
	expResp := exp.(*eerrs.Error)

	gotResp, exists := got.(*eerrs.Error)
	if !exists {
		return "no error occurred"
	}

	if gotResp.Code != expResp.Code {
		return "code does not match"
	}

	if gotResp.Message != expResp.Message {
		return "message does not match"
	}

	gotDetails := gotResp.Details.(errs.ExtraDetails)
	expDetails := expResp.Details.(errs.ExtraDetails)

	if gotDetails.HTTPStatus != expDetails.HTTPStatus {
		return "http status does not match"
	}

	if gotDetails.HTTPStatusCode != expDetails.HTTPStatusCode {
		return "http status code does not match"
	}

	return ""
}
