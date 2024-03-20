package mid

import (
	"context"

	"encore.dev/middleware"
)

type ctxKey int

const key ctxKey = 1

func setTraceID(req middleware.Request, traceID string) middleware.Request {
	ctx := context.WithValue(req.Context(), key, traceID)
	return req.WithContext(ctx)
}

// GetTraceID extracts the trace id for the request from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(string)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}

	return v
}

// =============================================================================

// Context adds values into the context of the request.
func Context(req middleware.Request, next middleware.Next) middleware.Response {
	req = setTraceID(req, req.Data().Trace.TraceID)

	return next(req)
}
