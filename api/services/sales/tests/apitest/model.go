package apitest

import (
	"context"

	eauth "encore.dev/beta/auth"
	"github.com/ardanlabs/encore/app/api/auth"
	"github.com/ardanlabs/encore/business/api/dbtest"
)

// User extends the dbtest user for app test support.
type User struct {
	dbtest.User
	Token string
}

// SeedData represents users for app tests.
type SeedData struct {
	Users  []User
	Admins []User
}

// Table represent fields needed for running an app test.
type Table struct {
	Name    string
	Token   string
	ExpResp any
	ExcFunc func(ctx context.Context) any
	CmpFunc func(got any, exp any) string
}

// AuthParams provides access to the authorization header.
type AuthParams struct {
	Authorization string `header:"Authorization"`
}

// AuthHandler represents a function that can perform authentication.
type AuthHandler func(ctx context.Context, ap *AuthParams) (eauth.UID, *auth.Claims, error)
