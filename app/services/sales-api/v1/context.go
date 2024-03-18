package encore

import (
	"context"
	"time"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/google/uuid"
)

type ctxKey int

const key ctxKey = 1

type values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

func getValues(ctx context.Context) *values {
	v, ok := ctx.Value(key).(*values)
	if !ok {
		return &values{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}

	return v
}

func getTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}

	return v.TraceID
}

func getTime(ctx context.Context) time.Time {
	v, ok := ctx.Value(key).(*values)
	if !ok {
		return time.Now()
	}

	return v.Now
}

func setValues(req middleware.Request, v *values) middleware.Request {
	ctx := context.WithValue(req.Context(), key, v)
	return req.WithContext(ctx)
}

// =============================================================================

type ctxUserKey int

const (
	userIDKey ctxUserKey = iota + 1
	userKey
)

func getUserID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}

	return v
}

func getUser(ctx context.Context) user.User {
	v, ok := ctx.Value(userKey).(user.User)
	if !ok {
		return user.User{}
	}

	return v
}

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func setUser(ctx context.Context, usr user.User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}
