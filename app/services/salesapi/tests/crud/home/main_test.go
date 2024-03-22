package home_test

import (
	"fmt"
	"os"
	"testing"

	"encore.dev/et"
	"github.com/ardanlabs/encore/app/services/salesapi"
	"github.com/ardanlabs/encore/business/data/dbtest"
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
		}

		t.Run(testName+"-"+tt.name, f)
	}
}
