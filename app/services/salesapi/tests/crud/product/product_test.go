package product_test

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

	// app.test(t, productQuery200(sd), "product-query-200")
	// app.test(t, productQueryByID200(sd), "product-querybyid-200")

	// app.test(t, productCreate200(sd), "product-create-200")
	// app.test(t, productCreate401(sd), "product-create-401")
	// app.test(t, productCreate400(sd), "product-create-400")

	// app.test(t, productUpdate200(sd), "product-update-200")
	// app.test(t, productUpdate401(sd), "product-update-401")
	// app.test(t, productUpdate400(sd), "product-update-400")

	// app.test(t, productDelete200(sd), "product-delete-200")
	// app.test(t, productDelete401(sd), "product-delete-401")
}
