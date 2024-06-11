package vproduct_test

import (
	"testing"
)

func Test_VProduct(t *testing.T) {
	t.Parallel()

	test := startTest(t)

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	test.Run(t, queryOk(sd), "query-ok")
}
