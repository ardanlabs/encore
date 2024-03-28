package dbtest

import (
	"context"
	"testing"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/business/api/errs"
)

// UnitTable represent fields needed for running an unit test.
type UnitTable struct {
	Name    string
	ExpResp any
	ExcFunc func(ctx context.Context) any
	CmpFunc func(got any, exp any) string
}

// =============================================================================

// UnitTest contains functions for executing an unit test.
type UnitTest struct{}

// Test performs the actual test logic based on the table data.
func (ut *UnitTest) Test(t *testing.T, table []UnitTable, testName string) {
	log := func(diff string, got any, exp any) {
		t.Log("DIFF")
		t.Logf("%s", diff)
		t.Log("GOT")
		t.Logf("%#v", got)
		t.Log("EXP")
		t.Logf("%#v", exp)
		t.Fatalf("Should get the expected response")
	}

	for _, tt := range table {
		f := func(t *testing.T) {
			ctx := context.Background()

			t.Log("Calling excFunc")
			got := tt.ExcFunc(ctx)

			diff := tt.CmpFunc(got, tt.ExpResp)
			if diff != "" {
				log(diff, got, tt.ExpResp)
			}
		}

		t.Run(testName+"-"+tt.Name, f)
	}
}

// =============================================================================

// CmpErrors compares two encore error values. If they are not equal, the
// reason is returned.
func CmpUnitErrors(got any, exp any) string {
	expResp := exp.(*eerrs.Error)

	gotResp, exists := got.(*eerrs.Error)
	if !exists {
		return "no error occurred"
	}

	if gotResp.Code != expResp.Code {
		return "code does not match"
	}

	if gotResp.Message != expResp.Message {
		return "message does not match"
	}

	gotDetails := gotResp.Details.(errs.ExtraDetails)
	expDetails := expResp.Details.(errs.ExtraDetails)

	if gotDetails.HTTPStatus != expDetails.HTTPStatus {
		return "http status does not match"
	}

	if gotDetails.HTTPStatusCode != expDetails.HTTPStatusCode {
		return "http status code does not match"
	}

	return ""
}
