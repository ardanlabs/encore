package mid

import (
	"context"
	"net/http"

	v1 "github.com/ardanlabs/encore/business/web/v1"
	"github.com/ardanlabs/encore/business/web/v1/auth"
	"github.com/ardanlabs/encore/foundation/logger"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/ardanlabs/encore/foundation/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) web.MidHandler {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Error(ctx, "message", "msg", err)

				var er v1.ErrorResponse
				var status int

				switch {
				case v1.IsTrustedError(err):
					trsErr := v1.GetTrustedError(err)

					if validate.IsFieldErrors(trsErr.Err) {
						fieldErrors := validate.GetFieldErrors(trsErr.Err)
						er = v1.ErrorResponse{
							Error:  "data validation error",
							Fields: fieldErrors.Fields(),
						}
						status = trsErr.Status
						break
					}

					er = v1.ErrorResponse{
						Error: trsErr.Error(),
					}
					status = trsErr.Status

				case auth.IsAuthError(err):
					er = v1.ErrorResponse{
						Error: http.StatusText(http.StatusUnauthorized),
					}
					status = http.StatusUnauthorized

				default:
					er = v1.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
