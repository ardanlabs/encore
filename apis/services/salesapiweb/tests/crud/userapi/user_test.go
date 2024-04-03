package user_test

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"encore.dev/et"
	"github.com/ardanlabs/encore/apis/services/salesapiweb"
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

	service, err := salesapiweb.NewService(dbTest.DB, dbTest.Auth)
	if err != nil {
		t.Fatalf("Service init error: %s", err)
	}
	et.MockService("salesapiweb", service, et.RunMiddleware(true))

	app := dbtest.AppTest{
		Service: service,
	}

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
