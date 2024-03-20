// Package appdb declares the application database, contains the SQL for
// database migrations and seeding. Provides a migration package for testing.
package appdb

import (
	edb "encore.dev/storage/sqldb"
)

// DBName represents the name of the database used by the application.
const DBName = "app"

// AppDB represents the database for this system. Encore will create and
// manage this database for us.
var AppDB = edb.NewDatabase("app", edb.DatabaseConfig{
	Migrations: "./migrate/migrations",
})
