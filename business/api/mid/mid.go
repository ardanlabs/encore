// Package mid provides middleware functions.
package mid

import (
	"context"
	"errors"

	eauth "encore.dev/beta/auth"
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

// GetUserID extracts the user id from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, found := eauth.UserID()
	if !found {
		return uuid.UUID{}, errors.New("user id not found")
	}

	v, err := uuid.Parse(string(userID))
	if err != nil {
		return uuid.UUID{}, err
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

func setProduct(req middleware.Request, prd product.Product) middleware.Request {
	ctx := context.WithValue(req.Context(), productKey, prd)
	return req.WithContext(ctx)
}

// GetProduct returns the product from the context.
func GetProduct(ctx context.Context) (product.Product, error) {
	v, ok := ctx.Value(productKey).(product.Product)
	if !ok {
		return product.Product{}, errors.New("product not found in context")
	}

	return v, nil
}

func setHome(req middleware.Request, hme home.Home) middleware.Request {
	ctx := context.WithValue(req.Context(), homeKey, hme)
	return req.WithContext(ctx)
}

// GetHome returns the home from the context.
func GetHome(ctx context.Context) (home.Home, error) {
	v, ok := ctx.Value(homeKey).(home.Home)
	if !ok {
		return home.Home{}, errors.New("home not found in context")
	}

	return v, nil
}
