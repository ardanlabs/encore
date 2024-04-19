// Package dbtest contains supporting code for running tests that hit the DB.
package dbtest

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/mail"
	"os/exec"
	"strings"
	"testing"
	"time"

	"encore.dev/rlog"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/api/delegate"
	"github.com/ardanlabs/encore/business/data/appdb/migrate"
	"github.com/ardanlabs/encore/business/data/sqldb"
	"github.com/ardanlabs/encore/business/domain/homebus"
	"github.com/ardanlabs/encore/business/domain/homebus/stores/homedb"
	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/ardanlabs/encore/business/domain/productbus/stores/productdb"
	"github.com/ardanlabs/encore/business/domain/userbus"
	"github.com/ardanlabs/encore/business/domain/userbus/stores/userdb"
	"github.com/ardanlabs/encore/business/domain/vproductbus"
	"github.com/ardanlabs/encore/business/domain/vproductbus/stores/vproductdb"
	"github.com/ardanlabs/encore/foundation/keystore"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
)

// Represents the secrets needed from the project.
var secrets struct {
	KeyID  string
	KeyPEM string
}

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

	fmt.Println("=== DB   ", url)

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
	Home     *homebus.Core
	Product  *productbus.Core
	User     *userbus.Core
	VProduct *vproductbus.Core
}

func newBusDomains(log rlog.Ctx, db *sqlx.DB) BusDomain {
	delegate := delegate.New(log)
	userBus := userbus.NewCore(log, delegate, userdb.NewStore(log, db))
	productBus := productbus.NewCore(log, userBus, delegate, productdb.NewStore(log, db))
	homeBus := homebus.NewCore(userBus, delegate, homedb.NewStore(log, db))
	vproductBus := vproductbus.NewCore(vproductdb.NewStore(log, db))

	return BusDomain{
		Delegate: delegate,
		Home:     homeBus,
		Product:  productBus,
		User:     userBus,
		VProduct: vproductBus,
	}
}

// =============================================================================

// Test owns state for running and shutting down tests.
type Test struct {
	Log       rlog.Ctx
	DB        *sqlx.DB
	Auth      *auth.Auth
	BusDomain BusDomain
	Teardown  func()
	t         *testing.T
}

// NewTest creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty. It returns
// the database to use as well as a function to call at the end of the test.
func NewTest(t *testing.T, url string, testName string) *Test {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	t.Logf("Creating database: %s", dbName)

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

	t.Logf("Migrating database: %s", dbName)

	if err := migrate.Migrate(ctx, db); err != nil {
		t.Fatalf("Migrating error: %s", err)
	}

	// -------------------------------------------------------------------------

	ks := keystore.New()
	if err := ks.LoadKey(secrets.KeyID, secrets.KeyPEM); err != nil {
		t.Fatalf("reading keys: %s", err)
	}

	a, err := auth.New(auth.Config{
		Log:       log,
		DB:        db,
		KeyLookup: ks,
	})
	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()

		db.Close()
		defer dbM.Close()

		t.Logf("Dropping database: %s", dbName)

		if _, err := dbM.ExecContext(context.Background(), "DROP DATABASE "+dbName); err != nil {
			fmt.Printf("dropping database %s: %v", dbName, err)
		}
	}

	tst := Test{
		Log:       log,
		DB:        db,
		Auth:      a,
		BusDomain: newBusDomains(log, db),
		Teardown:  teardown,
		t:         t,
	}

	return &tst
}

// Token generates an authenticated token for a user.
func (tst *Test) Token(email string) string {
	addr, _ := mail.ParseAddress(email)

	store := userdb.NewStore(tst.Log, tst.DB)
	dbUsr, err := store.QueryByEmail(context.Background(), *addr)
	if err != nil {
		return ""
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   dbUsr.ID.String(),
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: dbUsr.Roles,
	}

	token, err := tst.Auth.GenerateToken(secrets.KeyID, claims)
	if err != nil {
		tst.t.Fatal(err)
	}

	return token
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
