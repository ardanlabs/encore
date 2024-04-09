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

	dbTest := dbtest.NewTest(t, url, "Test_User")
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

	app.Test(t, userQueryOk(sd), "user-query-ok")
	app.Test(t, userQueryByIDOk(sd), "user-querybyid-ok")

	app.Test(t, userCreateOk(sd), "user-create-ok")
	app.Test(t, userCreateAuth(sd), "user-create-auth")
	app.Test(t, userCreateBad(sd), "user-create-bad")

	app.Test(t, userUpdateOk(sd), "user-update-ok")
	app.Test(t, userUpdateAuth(sd), "user-update-auth")
	app.Test(t, userUpdateBad(sd), "user-update-bad")

	app.Test(t, userDeleteOk(sd), "user-delete-ok")
	app.Test(t, userDeleteAuth(sd), "user-delete-auth")
}
