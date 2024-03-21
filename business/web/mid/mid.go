// Package mid provides middleware functions.
package mid

import (
	"context"
	"errors"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/google/uuid"
)

type ctxUserKey int
type ctxProductKey int

const (
	userIDKey  ctxUserKey    = 1
	userKey    ctxUserKey    = 2
	productKey ctxProductKey = 1
)

func setUserID(req middleware.Request, userID uuid.UUID) middleware.Request {
	ctx := context.WithValue(req.Context(), userIDKey, userID)
	return req.WithContext(ctx)
}

// GetUserID extracts the user id from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found")
	}

	return v, nil
}

func setUser(req middleware.Request, usr user.User) middleware.Request {
	ctx := context.WithValue(req.Context(), userKey, usr)
	return req.WithContext(ctx)
}

// GetUser extracts the user from the context.
func GetUser(ctx context.Context) (user.User, error) {
	v, ok := ctx.Value(userKey).(user.User)
	if !ok {
		return user.User{}, errors.New("user not found")
	}

	return v, nil
}

// GetProduct returns the product from the context.
func GetProduct(ctx context.Context) (product.Product, error) {
	v, ok := ctx.Value(productKey).(product.Product)
	if !ok {
		return product.Product{}, errors.New("product not found in context")
	}

	return v, nil
}

func setProduct(ctx context.Context, prd product.Product) context.Context {
	return context.WithValue(ctx, productKey, prd)
}
