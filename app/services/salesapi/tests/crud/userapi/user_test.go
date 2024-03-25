package user_test

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

	service, err := salesapi.NewService("user_test", dbTest.DB, dbTest.Auth)
	if err != nil {
		t.Fatalf("Service init error: %s", err)
	}
	et.MockService("salesapi", service, et.RunMiddleware(true))

	app := dbtest.AppTest{
		Service: service,
	}

	// -------------------------------------------------------------------------

	app.Test(t, userQuery200(sd), "user-query-200")
	app.Test(t, userQueryByID200(sd), "user-querybyid-200")

	app.Test(t, userCreate200(sd), "user-create-200")
	app.Test(t, userCreate401(sd), "user-create-401")
	app.Test(t, userCreate400(sd), "user-create-400")

	// app.test(t, userUpdate200(sd), "user-update-200")
	// app.test(t, userUpdate401(sd), "user-update-401")
	// app.test(t, userUpdate400(sd), "user-update-400")

	// app.test(t, userDelete200(sd), "user-delete-200")
	// app.test(t, userDelete401(sd), "user-delete-401")
}
