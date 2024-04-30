package product_test

import (
	"context"
	"testing"

	eauth "encore.dev/beta/auth"
	"encore.dev/et"
	authsrv "github.com/ardanlabs/encore/api/services/auth"
	salesrv "github.com/ardanlabs/encore/api/services/sales"
	"github.com/ardanlabs/encore/api/services/sales/tests/apitest"
	"github.com/ardanlabs/encore/app/api/auth"
	"github.com/ardanlabs/encore/app/api/mid"
	"github.com/ardanlabs/encore/business/api/dbtest"
)

func startTest(t *testing.T, url string, testName string) *apitest.Test {
	db := dbtest.NewDatabase(t, url, testName)

	// -------------------------------------------------------------------------

	ath, err := auth.New(auth.Config{
		Log:       db.Log,
		DB:        db.DB,
		KeyLookup: &apitest.KeyStore{},
	})
	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	authService, err := authsrv.NewService(db.Log, db.DB, ath)
	if err != nil {
		t.Fatalf("Auth service init error: %s", err)
	}
	et.MockService("auth", authService)

	salesService, err := salesrv.NewService(db.Log, db.DB)
	if err != nil {
		t.Fatalf("Sales service init error: %s", err)
	}
	et.MockService("sales", salesService, et.RunMiddleware(true))

	// -------------------------------------------------------------------------

	authHandler := func(ctx context.Context, ap *apitest.AuthParams) (eauth.UID, *auth.Claims, error) {
		return mid.Bearer(ctx, ath, ap.Authorization)
	}

	return apitest.New(db, ath, authHandler)
}
