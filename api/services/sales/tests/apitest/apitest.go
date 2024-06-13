// Package apitest contains supporting code for running app layer tests.
package apitest

import (
	"context"
	"net/mail"
	"testing"
	"time"

	eauth "encore.dev/beta/auth"
	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/sdk/auth"
	"github.com/ardanlabs/encore/business/domain/userbus/stores/userdb"
	"github.com/ardanlabs/encore/business/sdk/dbtest"
	"github.com/golang-jwt/jwt/v4"
)

// Test contains functions for executing an api test.
type Test struct {
	DB      *dbtest.Database
	Auth    *auth.Auth
	handler AuthHandler
}

// New constructs a Test value for running api tests.
func New(db *dbtest.Database, ath *auth.Auth, handler AuthHandler) *Test {
	return &Test{
		DB:      db,
		Auth:    ath,
		handler: handler,
	}
}

// Run performs the actual test logic based on the table data.
func (at *Test) Run(t *testing.T, table []Table, testName string) {
	log := func(diff string, got any, exp any) {

	}

	for _, tt := range table {
		f := func(t *testing.T) {
			t.Logf("\n***** Running Test: %s *****\n", testName+"-"+tt.Name)
			defer t.Logf("\n***** Finished Test: %s *****\n", testName+"-"+tt.Name)

			ctx := context.Background()

			ctx, err := at.authHandler(ctx, tt.Token)
			if err != nil {
				diff := tt.CmpFunc(err, tt.ExpResp)
				if diff != "" {
					log(diff, err, tt.ExpResp)
				}
				return
			}

			got := tt.ExcFunc(ctx)

			diff := tt.CmpFunc(got, tt.ExpResp)
			if diff != "" {
				t.Log("DIFF")
				t.Logf("%s", diff)
				t.Log("GOT")
				t.Logf("%#v", got)
				t.Log("EXP")
				t.Logf("%#v", tt.ExpResp)
				t.Fatalf("Should get the expected response")
			}
		}

		t.Run(testName+"-"+tt.Name, f)
	}
}

func (at *Test) authHandler(ctx context.Context, token string) (context.Context, error) {
	uid, claims, err := at.handler(ctx, &AuthParams{
		Authorization: "Bearer " + token,
	})

	if err != nil {
		return ctx, err
	}

	return eauth.WithContext(ctx, uid, claims), nil
}

// =============================================================================

// CmpAppErrors compares two encore error values. If they are not equal, the
// reason is returned.
func CmpAppErrors(got any, exp any) string {
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

	return ""
}

// Token generates an authenticated token for a user.
func Token(db *dbtest.Database, ath *auth.Auth, email string) string {
	addr, _ := mail.ParseAddress(email)

	store := userdb.NewStore(db.Log, db.DB)
	dbUsr, err := store.QueryByEmail(context.Background(), *addr)
	if err != nil {
		return ""
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   dbUsr.ID.String(),
			Issuer:    ath.Issuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: dbUsr.Roles,
	}

	token, err := ath.GenerateToken(kid, claims)
	if err != nil {
		return ""
	}

	return token
}
