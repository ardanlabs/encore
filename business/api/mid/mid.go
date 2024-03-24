// Package mid provides middleware functions.
package mid

import (
	"context"
	"errors"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/google/uuid"
)

type ctxKey int

const (
	userIDKey ctxKey = iota + 1
	userKey
	productKey
	homeKey
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

func setProduct(ctx context.Context, prd product.Product) context.Context {
	return context.WithValue(ctx, productKey, prd)
}

// GetProduct returns the product from the context.
func GetProduct(ctx context.Context) (product.Product, error) {
	v, ok := ctx.Value(productKey).(product.Product)
	if !ok {
		return product.Product{}, errors.New("product not found in context")
	}

	return v, nil
}

func setHome(ctx context.Context, hme home.Home) context.Context {
	return context.WithValue(ctx, homeKey, hme)
}

// GetHome returns the home from the context.
func GetHome(ctx context.Context) (home.Home, error) {
	v, ok := ctx.Value(homeKey).(home.Home)
	if !ok {
		return home.Home{}, errors.New("home not found in context")
	}

	return v, nil
}
