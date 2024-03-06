package sqldb

import (
	"context"
)

// NewDatabase declares a new SQL database.
//
// Encore uses static analysis to identify databases and their configuration,
// so all parameters passed to this function must be constant literals.
//
// A call to NewDatabase can only be made when declaring a package level variable. Any
// calls to this function made outside a package level variable declaration will result
// in a compiler error.
//
// The database name must be unique within the Encore application. Database names must be defined
// in kebab-case (lowercase alphanumerics and hyphen separated). Once created and deployed never
// change the database name, or else a new database will be created.
func NewDatabase(name string, config DatabaseConfig) (_ *Database) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L21-L23
	doPanic("encore apps must be run using the encore command")
	return
}

// DatabaseConfig specifies configuration for declaring a new database.
type DatabaseConfig struct {
	// Migrations is the directory containing the migration files
	// for this database.
	//
	// The path must be slash-separated relative path, and must be rooted within
	// the package directory (it cannot contain "../").
	// Valid paths are, for example, "migrations" or "db/migrations".
	//
	// Migrations are an ordered sequence of sql files of the format <number>_<description>.up.sql.
	Migrations string
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
//
// See (*database/sql.DB).ExecContext() for additional documentation.
func Exec(ctx context.Context, query string, args ...interface{}) (_ ExecResult, _ error) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L42-L44
	doPanic("encore apps must be run using the encore command")
	return
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
//
// See (*database/sql.DB).QueryContext() for additional documentation.
func Query(ctx context.Context, query string, args ...interface{}) (_ *Rows, _ error) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L50-L52
	doPanic("encore apps must be run using the encore command")
	return
}

// QueryRow executes a query that is expected to return at most one row.
//
// See (*database/sql.DB).QueryRowContext() for additional documentation.
func QueryRow(ctx context.Context, query string, args ...interface{}) (_ *Row) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L57-L59
	doPanic("encore apps must be run using the encore command")
	return
}

// Begin opens a new database transaction.
//
// See (*database/sql.DB).Begin() for additional documentation.
func Begin(ctx context.Context) (_ *Tx, _ error) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L64-L66
	doPanic("encore apps must be run using the encore command")
	return
}

// Commit commits the given transaction.
//
// See (*database/sql.Tx).Commit() for additional documentation.
// Deprecated: use tx.Commit() instead.
func Commit(tx *Tx) (_ error) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L72-L74
	doPanic("encore apps must be run using the encore command")
	return
}

// Rollback rolls back the given transaction.
//
// See (*database/sql.Tx).Rollback() for additional documentation.
// Deprecated: use tx.Rollback() instead.
func Rollback(tx *Tx) (_ error) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L80-L82
	doPanic("encore apps must be run using the encore command")
	return
}

// ExecTx is like Exec but executes the query in the given transaction.
//
// See (*database/sql.Tx).ExecContext() for additional documentation.
// Deprecated: use tx.Exec() instead.
func ExecTx(tx *Tx, ctx context.Context, query string, args ...interface{}) (_ ExecResult, _ error) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L88-L90
	doPanic("encore apps must be run using the encore command")
	return
}

// QueryTx is like Query but executes the query in the given transaction.
//
// See (*database/sql.Tx).QueryContext() for additional documentation.
// Deprecated: use tx.Query() instead.
func QueryTx(tx *Tx, ctx context.Context, query string, args ...interface{}) (_ *Rows, _ error) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L96-L98
	doPanic("encore apps must be run using the encore command")
	return
}

// QueryRowTx is like QueryRow but executes the query in the given transaction.
//
// See (*database/sql.Tx).QueryRowContext() for additional documentation.
// Deprecated: use tx.QueryRow() instead.
func QueryRowTx(tx *Tx, ctx context.Context, query string, args ...interface{}) (_ *Row) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L104-L106
	doPanic("encore apps must be run using the encore command")
	return
}

// constStr is a string that can only be provided as a constant.
type constStr string

// Named returns a database object connected to the database with the given name.
//
// The name must be a string literal constant, to facilitate static analysis.
func Named(name constStr) (_ *Database) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/pkgfn.go#L116-L118
	doPanic("encore apps must be run using the encore command")
	return
}
