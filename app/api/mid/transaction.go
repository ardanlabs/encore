package mid

import (
	"database/sql"
	"errors"

	eerrs "encore.dev/beta/errs"
	"encore.dev/middleware"
	"encore.dev/rlog"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/business/api/transaction"
)

// BeginCommitRollback starts a transaction for the domain call.
func BeginCommitRollback(log rlog.Ctx, bgn transaction.Beginner, req middleware.Request, next middleware.Next) middleware.Response {
	hasCommitted := false

	log.Info("BEGIN TRANSACTION")
	tx, err := bgn.Begin()
	if err != nil {
		return errs.NewResponsef(eerrs.Internal, "BEGIN TRANSACTION: %s", err)
	}

	defer func() {
		if !hasCommitted {
			log.Info("ROLLBACK TRANSACTION")
		}

		if err := tx.Rollback(); err != nil {
			if errors.Is(err, sql.ErrTxDone) {
				return
			}
			log.Info("ROLLBACK TRANSACTION", "ERROR", err)
		}
	}()

	req = setTran(req, tx)

	resp := next(req)
	if resp.Err != nil {
		return errs.NewResponsef(eerrs.Internal, "EXECUTE TRANSACTION: %s", resp.Err)
	}

	log.Info("COMMIT TRANSACTION")
	if err := tx.Commit(); err != nil {
		return errs.NewResponsef(eerrs.Internal, "COMMIT TRANSACTION: %s", err)
	}

	hasCommitted = true

	return resp
}
