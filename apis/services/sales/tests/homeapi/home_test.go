package home_test

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"encore.dev"
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
	if encore.Meta().Environment.Name == "ci-test" {
		return
	}

	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
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

func Test_Home(t *testing.T) {
	t.Parallel()

	dbTest, appTest := startTest(t, url, "Test_Home")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		dbTest.Teardown()
	}()

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(dbTest)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	appTest.Run(t, homeQueryOk(sd), "home-query-ok")
	appTest.Run(t, homeQueryByIDOk(sd), "home-querybyid-ok")

	appTest.Run(t, homeCreateOk(sd), "home-create-ok")
	appTest.Run(t, homeCreateBad(sd), "home-create-bad")
	appTest.Run(t, homeCreateAuth(sd), "home-create-auth")

	appTest.Run(t, homeUpdateOk(sd), "home-update-ok")
	appTest.Run(t, homeUpdateBad(sd), "home-update-bad")
	appTest.Run(t, homeUpdateAuth(sd), "home-update-auth")

	appTest.Run(t, homeDeleteOk(sd), "home-delete-ok")
	appTest.Run(t, homeDeleteAuth(sd), "home-delete-auth")
}

func startTest(t *testing.T, url string, testName string) (*dbtest.Test, *apptest.AppTest) {
	dbTest := dbtest.NewTest(t, url, testName)

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

	// -------------------------------------------------------------------------

	authHandler := func(ctx context.Context, ap *mid.AuthParams) (eauth.UID, *auth.Claims, error) {
		return mid.AuthHandler(ctx, dbTest.Auth, dbTest.BusDomain.User, ap)
	}

	appTest := apptest.New(authHandler)

	return dbTest, appTest
}
