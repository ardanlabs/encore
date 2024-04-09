package vproduct_test

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	eauth "encore.dev/beta/auth"
	"encore.dev/et"
	authsrv "github.com/ardanlabs/encore/apis/services/auth"
	"github.com/ardanlabs/encore/apis/services/sales"
	"github.com/ardanlabs/encore/app/api/apptest"
	"github.com/ardanlabs/encore/app/api/mid"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/data/dbtest"
)

var url string

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	et.EnableServiceInstanceIsolation()

	url, err = dbtest.StartDB()
	if err != nil {
		return 1, err
	}

	defer func() {
		err = dbtest.StopDB()
	}()

	return m.Run(), nil
}

// =============================================================================

func Test_VProduct(t *testing.T) {
	t.Parallel()

	dbTest := dbtest.NewTest(t, url, "Test_VProduct")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		dbTest.Teardown()
	}()

	sd, err := insertSeedData(dbTest)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	authService, err := authsrv.NewService(dbTest.Log, dbTest.DB, dbTest.Auth)
	if err != nil {
		t.Fatalf("Auth service init error: %s", err)
	}
	et.MockService("auth", authService)

	salesService, err := sales.NewService(dbTest.Log, dbTest.DB)
	if err != nil {
		t.Fatalf("Sales service init error: %s", err)
	}
	et.MockService("sales", salesService, et.RunMiddleware(true))

	authHandler := func(ctx context.Context, ap *mid.AuthParams) (eauth.UID, *auth.Claims, error) {
		return mid.AuthHandler(ctx, dbTest.Auth, dbTest.Core.BusCrud.User, ap)
	}

	app := apptest.New(authHandler)

	// -------------------------------------------------------------------------

	app.Test(t, vproductQueryOk(sd), "vproduct-query-ok")
}
