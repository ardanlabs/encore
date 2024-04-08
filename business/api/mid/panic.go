package mid

import (
	"runtime/debug"

	eerrs "encore.dev/beta/errs"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/app/api/metrics"
)

// Panics handles panics that occur when processing a request.
func Panics(v *metrics.Values, req middleware.Request, next middleware.Next) (resp middleware.Response) {
	defer func() {
		if rec := recover(); rec != nil {
			trace := debug.Stack()
			resp = errs.NewResponsef(eerrs.Internal, "PANIC [%v] TRACE[%s]", rec, string(trace))

			v.IncPanics()
		}
	}()

	return next(req)
}
