// Package migrate provides functions used by the dbtest package to init
// a database for proper use.
package migrate

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"sort"

	"github.com/ardanlabs/encore/business/data/sqldb"
	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*
var migrationDoc embed.FS

// Migrate will run the migration files.
func Migrate(ctx context.Context, db *sqlx.DB) (err error) {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	dirs, err := migrationDoc.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	files := make([]string, len(dirs))
	for i, entry := range dirs {
		files[i] = fmt.Sprintf("migrations/%s", entry.Name())
	}
	sort.Strings(files)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if errTx := tx.Rollback(); errTx != nil {
			if errors.Is(errTx, sql.ErrTxDone) {
				return
			}

			err = fmt.Errorf("rollback: %w", errTx)
			return
		}
	}()

	for _, file := range files {
		doc, err := migrationDoc.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read file: %s: %w", file, err)
		}

		if _, err := tx.Exec(string(doc)); err != nil {
			return fmt.Errorf("exec: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

//go:embed seeds/seed.sql
var seedDoc string

// Seed will insert data needed for a new database.
func Seed(ctx context.Context, db *sqlx.DB) (err error) {
	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if errTx := tx.Rollback(); errTx != nil {
			if errors.Is(errTx, sql.ErrTxDone) {
				return
			}

			err = fmt.Errorf("rollback: %w", errTx)
			return
		}
	}()

	if _, err := tx.Exec(seedDoc); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
