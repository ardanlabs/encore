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

	test.Run(t, productQueryOk(sd), "product-query-ok")
	test.Run(t, productQueryByIDOk(sd), "product-querybyid-ok")

	test.Run(t, productCreateOk(sd), "product-create-ok")
	test.Run(t, productCreateBad(sd), "product-create-bad")
	test.Run(t, productCreateAuth(sd), "product-create-auth")

	test.Run(t, productUpdateOk(sd), "product-update-ok")
	test.Run(t, productUpdateBad(sd), "product-update-bad")
	test.Run(t, productUpdateAuth(sd), "product-update-auth")

	test.Run(t, productDeleteOk(sd), "product-delete-ok")
	test.Run(t, productDeleteAuth(sd), "product-delete-auth")
}
