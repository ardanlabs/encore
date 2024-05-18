// Package dbtest contains supporting code for running tests that hit the DB.
package dbtest

import (
	"context"
	"testing"
	"time"

	esqldb "encore.dev/storage/sqldb"
	"github.com/ardanlabs/encore/business/domain/homebus"
	"github.com/ardanlabs/encore/business/domain/homebus/stores/homedb"
	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/ardanlabs/encore/business/domain/productbus/stores/productdb"
	"github.com/ardanlabs/encore/business/domain/userbus"
	"github.com/ardanlabs/encore/business/domain/userbus/stores/usercache"
	"github.com/ardanlabs/encore/business/domain/userbus/stores/userdb"
	"github.com/ardanlabs/encore/business/domain/vproductbus"
	"github.com/ardanlabs/encore/business/domain/vproductbus/stores/vproductdb"
	"github.com/ardanlabs/encore/business/sdk/delegate"
	"github.com/ardanlabs/encore/business/sdk/sqldb"
	"github.com/ardanlabs/encore/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	Delegate *delegate.Delegate
	Home     *homebus.Business
	Product  *productbus.Business
	User     *userbus.Business
	VProduct *vproductbus.Business
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {
	delegate := delegate.New(log)
	userBus := userbus.NewBusiness(log, delegate, usercache.NewStore(log, userdb.NewStore(log, db), time.Hour))
	productBus := productbus.NewBusiness(log, userBus, delegate, productdb.NewStore(log, db))
	homeBus := homebus.NewBusiness(log, userBus, delegate, homedb.NewStore(log, db))
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
	DB        *sqlx.DB
	Log       *logger.Logger
	BusDomain BusDomain
	Teardown  func()
}

// NewDatabase uses the specified database to perform testing. This database
// should be created using the `et.NewTestDatabase` call. A connection pool
// is provided with business domain packages.
func NewDatabase(t *testing.T, edb *esqldb.Database) *Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := sqldb.Open(sqldb.Config{
		EDB:          edb,
		MaxIdleConns: 0,
		MaxOpenConns: 0,
	})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	if err := sqldb.StatusCheck(ctx, db); err != nil {
		t.Fatalf("status check database: %v", err)
	}

	// -------------------------------------------------------------------------

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
	}

	log := logger.New("test")

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

// UserNamePointer is a helper to get a *Name from a string. It's in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func UserNamePointer(value string) *userbus.Name {
	name := userbus.Names.MustParse(value)
	return &name
}

// ProductNamePointer is a helper to get a *Name from a string. It's in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func ProductNamePointer(value string) *productbus.Name {
	name := productbus.Names.MustParse(value)
	return &name
}
