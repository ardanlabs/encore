package mid

import (
	"net/http"

	"encore.dev/middleware"
	v1 "github.com/ardanlabs/encore/business/web/v1"
	"github.com/ardanlabs/encore/business/web/v1/auth"
	"github.com/ardanlabs/encore/foundation/logger"
	"github.com/ardanlabs/encore/foundation/validate"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger, req middleware.Request, next middleware.Next) middleware.Response {
	resp := next(req)
	if resp.Err == nil {
		return resp
	}

	ctx := req.Context()

	log.Error(ctx, "errors message", "msg", resp.Err)

	var midResp middleware.Response

	switch {
	case v1.IsTrustedError(resp.Err):
		trsErr := v1.GetTrustedError(resp.Err)

		if validate.IsFieldErrors(trsErr.Err) {
			fieldErrors := validate.GetFieldErrors(trsErr.Err)
			midResp = v1.NewErrorResponseWithFields(trsErr.Status, "data validation error", fieldErrors)
			break
		}

		midResp = v1.NewErrorResponse(trsErr.Status, trsErr)

	case auth.IsAuthError(resp.Err):
		midResp = v1.NewErrorResponseWithMessage(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))

	default:
		midResp = v1.NewErrorResponseWithMessage(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return midResp
}
