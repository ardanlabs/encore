package tests

import (
	"runtime/debug"
	"testing"

	"encore.dev/et"
	"github.com/ardanlabs/encore/app/services/sales-api/encore"
	"github.com/ardanlabs/encore/business/data/dbtest"
)

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

	service, err := encore.InitService(dbTest.DB, "../../../../../zarf/keys")
	if err != nil {
		t.Fatalf("Service init error: %s", err)
	}
	et.MockService("encore", service)

	app := appTest{
		service:    service,
		userToken:  dbTest.TokenV1("user@example.com", "gophers"),
		adminToken: dbTest.TokenV1("admin@example.com", "gophers"),
	}

	// -------------------------------------------------------------------------

	sd, err := createUserSeed(dbTest)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	app.test(t, userQuery200(sd), "user-query-200")
	// app.test(t, userQueryByID200(sd), "user-querybyid-200")

	// app.test(t, userCreate200(sd), "user-create-200")
	// app.test(t, userCreate401(sd), "user-create-401")
	// app.test(t, userCreate400(sd), "user-create-400")

	// app.test(t, userUpdate200(sd), "user-update-200")
	// app.test(t, userUpdate401(sd), "user-update-401")
	// app.test(t, userUpdate400(sd), "user-update-400")

	// app.test(t, userDelete200(sd), "user-delete-200")
	// app.test(t, userDelete401(sd), "user-delete-401")
}
