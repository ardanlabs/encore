package home_test

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"encore.dev/et"
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

	dbTest := dbtest.NewTest(t, url, "Test_Home/crud")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		dbTest.Teardown()
	}()

	// sd, err := insertSeedData(dbTest)
	// if err != nil {
	// 	t.Fatalf("Seeding error: %s", err)
	// }

	// -------------------------------------------------------------------------

	// service, err := salesapi.NewService(dbTest.DB, dbTest.Auth)
	// if err != nil {
	// 	t.Fatalf("Service init error: %s", err)
	// }
	// et.MockService("salesapi", service)

	// app := dbtest.AppTest{
	// 	Service: service,
	// }

	// -------------------------------------------------------------------------

	// app.test(t, homeQuery200(sd), "home-query-200")
	// app.test(t, homeQueryByID200(sd), "home-querybyid-200")

	// app.test(t, homeCreate200(sd), "home-create-200")
	// app.test(t, homeCreate401(sd), "home-create-401")
	// app.test(t, homeCreate400(sd), "home-create-400")

	// app.test(t, homeUpdate200(sd), "home-update-200")
	// app.test(t, homeUpdate401(sd), "home-update-401")
	// app.test(t, homeUpdate400(sd), "home-update-400")

	// app.test(t, homeDelete200(sd), "home-delete-200")
	// app.test(t, homeDelete401(sd), "home-delete-401")
}
