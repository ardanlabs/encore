package home_test

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	eauth "encore.dev/beta/auth"
	"encore.dev/et"
	"github.com/ardanlabs/encore/apis/sales"
	"github.com/ardanlabs/encore/app/api/apptest"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/api/mid"
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

func Test_Home(t *testing.T) {
	t.Parallel()

	dbTest := dbtest.NewTest(t, url, "Test_Home")
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

	service, err := sales.NewService(dbTest.Log, dbTest.DB)
	if err != nil {
		t.Fatalf("Service init error: %s", err)
	}
	et.MockService("sales", service, et.RunMiddleware(true))

	authHandler := func(ctx context.Context, ap *mid.AuthParams) (eauth.UID, *auth.Claims, error) {
		return "", nil, nil
	}

	app := apptest.New(authHandler)

	// -------------------------------------------------------------------------

	app.Test(t, homeQueryOk(sd), "home-query-ok")
	app.Test(t, homeQueryByIDOk(sd), "home-querybyid-ok")

	app.Test(t, homeCreateOk(sd), "home-create-ok")
	app.Test(t, homeCreateBad(sd), "home-create-bad")
	app.Test(t, homeCreateAuth(sd), "home-create-auth")

	app.Test(t, homeUpdateOk(sd), "home-update-ok")
	app.Test(t, homeUpdateBad(sd), "home-update-bad")
	app.Test(t, homeUpdateAuth(sd), "home-update-auth")

	app.Test(t, homeDeleteOk(sd), "home-delete-ok")
	app.Test(t, homeDeleteAuth(sd), "home-delete-auth")
}
