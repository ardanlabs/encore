package database

import (
	edb "encore.dev/storage/sqldb"
)

const DBName = "app"

// We are declaring the existence of a database for this system. It MUST
// be declared at a package level with the Service type.
var EDB = edb.NewDatabase("app", edb.DatabaseConfig{
	Migrations: "./migrate/migrations",
})
