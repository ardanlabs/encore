package home_test

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"encore.dev/et"
	"github.com/ardanlabs/encore/app/services/salesapi"
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

	service, err := salesapi.NewService(dbTest.DB, dbTest.Auth)
	if err != nil {
		t.Fatalf("Service init error: %s", err)
	}
	et.MockService("salesapi", service, et.RunMiddleware(true))

	app := dbtest.AppTest{
		Service: service,
	}

	// -------------------------------------------------------------------------

	app.Test(t, homeQueryOk(sd), "home-query-ok")
	// app.Test(t, homeQueryByIDOk(sd), "home-querybyid-ok")

	// app.Test(t, homeCreateOk(sd), "home-create-ok")
	// app.Test(t, homeCreateBad(sd), "home-create-bad")
	// app.Test(t, homeCreateAuth(sd), "home-create-auth")

	// app.Test(t, homeUpdateOk(sd), "home-update-ok")
	// app.Test(t, homeUpdateBad(sd), "home-update-bad")
	// app.Test(t, homeUpdateAuth(sd), "home-update-auth")

	// app.Test(t, homeDeleteOk(sd), "home-delete-ok")
	// app.Test(t, homeDeleteAuth(sd), "home-delete-auth")
}
