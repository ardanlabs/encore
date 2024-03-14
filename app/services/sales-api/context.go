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

// Values represent state for each request.
type values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// GetValues returns the values from the context.
func GetValues(ctx context.Context) *values {
	v, ok := ctx.Value(key).(*values)
	if !ok {
		return &values{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}

	return v
}

// GetTraceID returns the trace id from the context.
func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value(key).(*values)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}

	return v.TraceID
}

// GetTime returns the time from the context.
func GetTime(ctx context.Context) time.Time {
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

// GetUserID returns the claims from the context.
func GetUserID(ctx context.Context) uuid.UUID {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}
	}

	return v
}

// GetUser returns the user from the context.
func GetUser(ctx context.Context) user.User {
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
