package home_test

import (
	"testing"
)

func Test_Home(t *testing.T) {
	t.Parallel()

	test := startTest(t)

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
