package encore

import (
	"context"

	"encore.dev/middleware"
)

type ctxKey int

const key ctxKey = 1

type values struct {
	TraceID string
}

func setValues(req middleware.Request, v *values) middleware.Request {
	ctx := context.WithValue(req.Context(), key, v)
	return req.WithContext(ctx)
}

func getTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}

	return v.TraceID
}
