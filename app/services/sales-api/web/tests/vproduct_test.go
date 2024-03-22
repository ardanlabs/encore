package tests

import (
	"runtime/debug"
	"testing"

	"encore.dev/et"
	"github.com/ardanlabs/encore/app/services/sales-api/encore"
	"github.com/ardanlabs/encore/business/data/dbtest"
)

func Test_VProduct(t *testing.T) {
	t.Parallel()

	dbTest := dbtest.NewTest(t, url, "Test_VProduct")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		dbTest.Teardown()
	}()

	service, err := encore.InitService(dbTest.DB, "../../../../../zarf/keys")
	if err != nil {
		t.Fatalf("Service init error: %s", err)
	}
	et.MockService("encore", service)

	// -------------------------------------------------------------------------

	// sd, err := createProductSeed(dbTest)
	// if err != nil {
	// 	t.Fatalf("Seeding error: %s", err)
	// }

	// -------------------------------------------------------------------------

	//app.test(t, vproductQuery200(sd), "vproduct-query-200")
}
