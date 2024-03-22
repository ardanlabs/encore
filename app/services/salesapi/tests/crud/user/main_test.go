package user_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"encore.dev/beta/auth"
	"encore.dev/et"
	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/business/data/dbtest"
	"github.com/ardanlabs/encore/business/web/mid"
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

	fmt.Println("URL:", url)

	return m.Run(), nil
}

type appTest struct {
	service    *salesapi.Service
	userToken  string
	adminToken string
}

func (at *appTest) test(t *testing.T, table []tableData, testName string) {
	for _, tt := range table {
		f := func(t *testing.T) {
			resp := tt.excFunc(at.service)

			diff := tt.cmpFunc(resp, tt.expResp)
			if diff != "" {
				t.Log("GOT")
				t.Logf("%#v", tt.resp)
				t.Log("EXP")
				t.Logf("%#v", tt.expResp)
				t.Fatalf("Should get the expected response")
			}
		}

		t.Run(testName+"-"+tt.name, f)
	}
}

func authHandler(ctx context.Context, s *salesapi.Service, token string) (context.Context, error) {
	uid, claims, err := s.AuthHandler(ctx, &mid.AuthParams{
		Authorization: "Bearer " + token,
	})

	if err != nil {
		return ctx, err
	}

	return auth.WithContext(ctx, uid, claims), nil
}
