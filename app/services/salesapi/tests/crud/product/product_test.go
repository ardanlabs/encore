package product_test

import (
	"runtime/debug"
	"testing"

	"encore.dev/et"
	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/business/data/dbtest"
)

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

	service, err := salesapi.NewService(dbTest.DB, dbTest.Auth)
	if err != nil {
		t.Fatalf("Service init error: %s", err)
	}
	et.MockService("salesapi", service)

	// -------------------------------------------------------------------------

	// sd, err := createProductSeed(dbTest)
	// if err != nil {
	// 	t.Fatalf("Seeding error: %s", err)
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
