// Package appdb declares the application database, contains the SQL for
// database migrations and seeding. Provides a migration package for testing.
package appdb

import (
	edb "encore.dev/storage/sqldb"
)

// This represents the database for this system. Encore will create and
// manage this database for us. The name has to be a literal string.
var _ = edb.NewDatabase("app", edb.DatabaseConfig{
	Migrations: "./migrate/migrations",
})
