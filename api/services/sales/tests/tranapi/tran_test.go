package tran_test

import (
	"runtime/debug"
	"testing"
)

func Test_Tran(t *testing.T) {
	t.Parallel()

	test := startTest(t)
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

	test.Run(t, createOk(sd), "create-ok")
}
