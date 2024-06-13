// Package migrate provides functions used by the dbtest package to init
// a database for proper use.
package migrate

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/ardanlabs/encore/business/sdk/sqldb"
	"github.com/jmoiron/sqlx"
)

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
