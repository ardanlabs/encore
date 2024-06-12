package mid

import (
	"context"
	"database/sql"
	"errors"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/sdk/errs"
	"github.com/ardanlabs/encore/business/sdk/sqldb"
	"github.com/ardanlabs/encore/foundation/logger"
)

// BeginCommitRollback starts a transaction for the domain call.
func BeginCommitRollback(log *logger.Logger, bgn sqldb.Beginner, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := context.Background()

	hasCommitted := false

	log.Info(ctx, "BEGIN TRANSACTION")
	tx, err := bgn.Begin()
	if err != nil {
		return errs.NewResponsef(errs.Internal, "BEGIN TRANSACTION: %s", err)
	}

	defer func() {
		if !hasCommitted {
			log.Info(ctx, "ROLLBACK TRANSACTION")
		}

		if err := tx.Rollback(); err != nil {
			if errors.Is(err, sql.ErrTxDone) {
				return
			}
			log.Info(ctx, "ROLLBACK TRANSACTION", "ERROR", err)
		}
	}()

	req = setTran(req, tx)

	resp := next(req)
	if resp.Err != nil {
		return errs.NewResponsef(errs.Internal, "EXECUTE TRANSACTION: %s", resp.Err)
	}

	log.Info(ctx, "COMMIT TRANSACTION")
	if err := tx.Commit(); err != nil {
		return errs.NewResponsef(errs.Internal, "COMMIT TRANSACTION: %s", err)
	}

	hasCommitted = true

	return resp
}
