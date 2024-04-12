package user_test

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
	et.EnableServiceInstanceIsolation()

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

func Test_User(t *testing.T) {
	t.Parallel()

	dbTest, appTest := startTest(t, url, "Test_User")
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

	appTest.Run(t, userQueryOk(sd), "user-query-ok")
	appTest.Run(t, userQueryByIDOk(sd), "user-querybyid-ok")

	appTest.Run(t, userCreateOk(sd), "user-create-ok")
	appTest.Run(t, userCreateAuth(sd), "user-create-auth")
	appTest.Run(t, userCreateBad(sd), "user-create-bad")

	appTest.Run(t, userUpdateOk(sd), "user-update-ok")
	appTest.Run(t, userUpdateAuth(sd), "user-update-auth")
	appTest.Run(t, userUpdateBad(sd), "user-update-bad")

	appTest.Run(t, userDeleteOk(sd), "user-delete-ok")
	appTest.Run(t, userDeleteAuth(sd), "user-delete-auth")
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
