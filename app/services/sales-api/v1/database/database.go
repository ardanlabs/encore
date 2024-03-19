package database

import (
	_ "embed"

	edb "encore.dev/storage/sqldb"
)

// We are declaring the existence of a database for this system. It MUST
// be declared at a package level with the Service type.
var EDB = edb.NewDatabase("app", edb.DatabaseConfig{
	Migrations: "./migrate/migrations",
})
