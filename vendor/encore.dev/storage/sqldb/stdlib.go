package sqldb

// RegisterStdlibDriver returns a connection string that can be used with
// the standard library's sql.Open function to connect to the same db.
//
// The connection string should be used with the "encore" driver name:
//
//	connStr := sqldb.RegisterStdlibDriver(myDB)
//	db, err := sql.Open("encore", connStr)
//
// The main use case is to support libraries that expect to call sql.Open
// themselves without exposing the underlying database credentials.
func RegisterStdlibDriver(db *Database) (_ string) {
	// Encore will provide an implementation to this function at runtime, we do not expose
	// the implementation in the API contract as it is an implementation detail, which may change
	// between releases.
	//
	// The current implementation of this function can be found here:
	//    https://github.com/encoredev/encore/blob/v1.30.0/runtimes/go/storage/sqldb/stdlib.go#L22-L31
	doPanic("encore apps must be run using the encore command")
	return
}
