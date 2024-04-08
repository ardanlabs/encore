package product_test

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"encore.dev/et"
	"github.com/ardanlabs/encore/apis/services/salesapiweb"
	"github.com/ardanlabs/encore/app/api/apptest"
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

func Test_Product(t *testing.T) {
	t.Parallel()

	dbTest := dbtest.NewTest(t, url, "Test_Product")
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

	service, err := salesapiweb.NewService(dbTest.Log, dbTest.DB, dbTest.Auth)
	if err != nil {
		t.Fatalf("Service init error: %s", err)
	}
	et.MockService("salesapiweb", service, et.RunMiddleware(true))

	app := apptest.New(service)

	// -------------------------------------------------------------------------

	app.Test(t, productQueryOk(sd), "product-query-ok")
	app.Test(t, productQueryByIDOk(sd), "product-querybyid-ok")

	app.Test(t, productCreateOk(sd), "product-create-ok")
	app.Test(t, productCreateBad(sd), "product-create-bad")
	app.Test(t, productCreateAuth(sd), "product-create-auth")

	app.Test(t, productUpdateOk(sd), "product-update-ok")
	app.Test(t, productUpdateBad(sd), "product-update-bad")
	app.Test(t, productUpdateAuth(sd), "product-update-auth")

	app.Test(t, productDeleteOk(sd), "product-delete-ok")
	app.Test(t, productDeleteAuth(sd), "product-delete-auth")
}
