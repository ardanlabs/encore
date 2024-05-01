package product_test

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

func Test_Product(t *testing.T) {
	t.Parallel()

	test := startTest(t, url, "Test_Product")
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

	test.Run(t, queryOk(sd), "query-ok")
	test.Run(t, queryByIDOk(sd), "querybyid-ok")

	test.Run(t, createOk(sd), "create-ok")
	test.Run(t, createBad(sd), "create-bad")
	test.Run(t, createAuth(sd), "create-auth")

	test.Run(t, updateOk(sd), "update-ok")
	test.Run(t, updateBad(sd), "update-bad")
	test.Run(t, updateAuth(sd), "update-auth")

	test.Run(t, deleteOk(sd), "delete-ok")
	test.Run(t, deleteAuth(sd), "delete-auth")
}
