package encore

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

func getTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(string)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}

	return v
}
