package home_test

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"encore.dev"
	"github.com/ardanlabs/encore/business/api/dbtest"
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

	test := startTest(t, url, "Test_Home")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		test.DB.Teardown()
	}()

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	test.Run(t, homeQueryOk(sd), "home-query-ok")
	test.Run(t, homeQueryByIDOk(sd), "home-querybyid-ok")

	test.Run(t, homeCreateOk(sd), "home-create-ok")
	test.Run(t, homeCreateBad(sd), "home-create-bad")
	test.Run(t, homeCreateAuth(sd), "home-create-auth")

	test.Run(t, homeUpdateOk(sd), "home-update-ok")
	test.Run(t, homeUpdateBad(sd), "home-update-bad")
	test.Run(t, homeUpdateAuth(sd), "home-update-auth")

	test.Run(t, homeDeleteOk(sd), "home-delete-ok")
	test.Run(t, homeDeleteAuth(sd), "home-delete-auth")
}
