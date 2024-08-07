package tran_test

import (
	"testing"
)

func Test_Tran(t *testing.T) {
	t.Parallel()

	test := startTest(t)

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	test.Run(t, createOk(sd), "create-ok")
}
