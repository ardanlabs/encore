// Package dbtest contains supporting code for running tests that hit the DB.
package dbtest

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"testing"
	"time"

	"encore.dev/rlog"
	"github.com/ardanlabs/encore/business/api/appdb/migrate"
	"github.com/ardanlabs/encore/business/api/delegate"
	"github.com/ardanlabs/encore/business/api/sqldb"
	"github.com/ardanlabs/encore/business/domain/homebus"
	"github.com/ardanlabs/encore/business/domain/homebus/stores/homedb"
	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/ardanlabs/encore/business/domain/productbus/stores/productdb"
	"github.com/ardanlabs/encore/business/domain/userbus"
	"github.com/ardanlabs/encore/business/domain/userbus/stores/usercache"
	"github.com/ardanlabs/encore/business/domain/userbus/stores/userdb"
	"github.com/ardanlabs/encore/business/domain/vproductbus"
	"github.com/ardanlabs/encore/business/domain/vproductbus/stores/vproductdb"
	"github.com/jmoiron/sqlx"
)

// StartDB retrieves the database information.
func StartDB() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("encore", "db", "conn-uri", "--test", "app")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("could not access the database information: %w", err)
	}

	url := out.String()
	url = strings.Trim(url, "\n")

	return url, nil
}

// StopDB stops a running database instance.
func StopDB() error {
	return nil
}

// =============================================================================

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	Delegate *delegate.Delegate
	Home     *homebus.Business
	Product  *productbus.Business
	User     *userbus.Business
	VProduct *vproductbus.Business
}

func newBusDomains(log rlog.Ctx, db *sqlx.DB) BusDomain {
	delegate := delegate.New(log)
	userBus := userbus.NewBusiness(log, delegate, usercache.NewStore(log, userdb.NewStore(log, db), time.Hour))
	productBus := productbus.NewBusiness(log, userBus, delegate, productdb.NewStore(log, db))
	homeBus := homebus.NewBusiness(userBus, delegate, homedb.NewStore(log, db))
	vproductBus := vproductbus.NewBusiness(vproductdb.NewStore(log, db))

	return BusDomain{
		Delegate: delegate,
		Home:     homeBus,
		Product:  productBus,
		User:     userBus,
		VProduct: vproductBus,
	}
}

// =============================================================================

// Database owns state for running and shutting down tests.
type Database struct {
	Log       rlog.Ctx
	DB        *sqlx.DB
	BusDomain BusDomain
	Teardown  func()
}

// NewDatabase creates a new test database inside the database that was started
// to handle testing. The database is migrated to the current version and
// a connection pool is provided with business domain packages.
func NewDatabase(t *testing.T, url string, testName string) *Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Logf("DB: %s\n", url)

	log := rlog.With("service", "sales-test")

	dbM, err := sqldb.OpenTest(url)
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	if err := sqldb.StatusCheck(ctx, dbM); err != nil {
		t.Fatalf("status check database: %v", err)
	}

	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 4)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	dbName := string(b)

	t.Logf("Creating Database: %s", dbName)
	if _, err := dbM.ExecContext(context.Background(), "CREATE DATABASE "+dbName); err != nil {
		t.Fatalf("creating database %s: %v", dbName, err)
	}

	// -------------------------------------------------------------------------

	// This is changing out the base dbname with the new one on
	// the connection string.
	url = strings.Replace(url, "app", dbName, 1)

	db, err := sqldb.OpenTest(url)
	if err != nil {
		t.Fatalf("Opening database connection: %v", err)
	}

	t.Logf("Migrating Database: %s", dbName)
	if err := migrate.Migrate(ctx, db); err != nil {
		t.Fatalf("Migrating error: %s", err)
	}

	// -------------------------------------------------------------------------

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()

		db.Close()
		defer dbM.Close()

		t.Logf("Dropping Database: %s", dbName)
		if _, err := dbM.ExecContext(context.Background(), "DROP DATABASE "+dbName); err != nil {
			fmt.Printf("dropping database %s: %v", dbName, err)
		}
	}

	return &Database{
		Log:       log,
		DB:        db,
		BusDomain: newBusDomains(log, db),
		Teardown:  teardown,
	}
}

// =============================================================================

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}

// FloatPointer is a helper to get a *float64 from a float64. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func FloatPointer(f float64) *float64 {
	return &f
}

// BoolPointer is a helper to get a *bool from a bool. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func BoolPointer(b bool) *bool {
	return &b
}
