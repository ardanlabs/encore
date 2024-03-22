package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/encore/app/services/sales-api/encore"
	"github.com/ardanlabs/encore/business/data/dbtest"
)

var url string
var service *encore.Service

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
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
	service    *encore.Service
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
