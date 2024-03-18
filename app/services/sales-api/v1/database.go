package encore

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	edb "encore.dev/storage/sqldb"
	"github.com/ardanlabs/encore/business/data/sqldb"
	"github.com/jmoiron/sqlx"
)

// We are declaring the existence of a database for this system. It MUST
// be declared at a package level with the Service type.
var ebdDB = edb.NewDatabase("url", edb.DatabaseConfig{
	Migrations: "./migrations",
})

//go:embed seeds/seed.sql
var seedDoc string

func seedDatabase(ctx context.Context, db *sqlx.DB) (err error) {
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
