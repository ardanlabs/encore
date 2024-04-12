// Package mid provides context support.
package mid

import (
	"context"
	"errors"

	eauth "encore.dev/beta/auth"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/domain/homebus"
	"github.com/ardanlabs/encore/business/domain/productbus"
	"github.com/ardanlabs/encore/business/domain/userbus"
	"github.com/google/uuid"
)

// AuthInfo defines the information required to perform an authorization.
type AuthInfo struct {
	Claims auth.Claims
	UserID uuid.UUID
	Rule   string
}

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

func setUser(req middleware.Request, usr userbus.User) middleware.Request {
	ctx := context.WithValue(req.Context(), userKey, usr)
	return req.WithContext(ctx)
}

// GetUser extracts the user from the context.
func GetUser(ctx context.Context) (userbus.User, error) {
	v, ok := ctx.Value(userKey).(userbus.User)
	if !ok {
		return userbus.User{}, errors.New("user not found")
	}

	return v, nil
}

func setProduct(req middleware.Request, prd productbus.Product) middleware.Request {
	ctx := context.WithValue(req.Context(), productKey, prd)
	return req.WithContext(ctx)
}

// GetProduct returns the product from the context.
func GetProduct(ctx context.Context) (productbus.Product, error) {
	v, ok := ctx.Value(productKey).(productbus.Product)
	if !ok {
		return productbus.Product{}, errors.New("product not found in context")
	}

	return v, nil
}

func setHome(req middleware.Request, hme homebus.Home) middleware.Request {
	ctx := context.WithValue(req.Context(), homeKey, hme)
	return req.WithContext(ctx)
}

// GetHome returns the home from the context.
func GetHome(ctx context.Context) (homebus.Home, error) {
	v, ok := ctx.Value(homeKey).(homebus.Home)
	if !ok {
		return homebus.Home{}, errors.New("home not found in context")
	}

	return v, nil
}
