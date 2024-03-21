package tests

import (
	"runtime/debug"
	"testing"

	"github.com/ardanlabs/encore/business/data/dbtest"
)

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

	// userToken := test.TokenV1("user@example.com", "gophers")
	// adminToken := test.TokenV1("admin@example.com", "gophers")

	// -------------------------------------------------------------------------

	// sd, err := createHomeSeed(dbTest)
	// if err != nil {
	// 	t.Fatalf("Seeding error: %s", err)
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
