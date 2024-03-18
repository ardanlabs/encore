// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ardanlabs/encore/business/core/crud/user"
	v1 "github.com/ardanlabs/encore/business/web/v1"
	"github.com/ardanlabs/encore/business/web/v1/auth"
)

type Handlers struct {
	user *user.Core
	auth *auth.Auth
}

func New(user *user.Core, auth *auth.Auth) *Handlers {
	return &Handlers{
		user: user,
		auth: auth,
	}
}

// Create adds a new user to the system.
func (h *Handlers) Create(ctx context.Context, app AppNewUser) (AppUser, error) {
	nc, err := toCoreNewUser(app)
	if err != nil {
		return AppUser{}, v1.NewTrustedError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return AppUser{}, v1.NewTrustedError(user.ErrUniqueEmail, http.StatusConflict)
		}
		return AppUser{}, fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return toAppUser(usr), nil
}

// Token provides an API token for the authenticated user.
// func (h *Handlers) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
// 	kid := web.Param(r, "kid")
// 	if kid == "" {
// 		return validate.NewFieldsError("kid", errors.New("missing kid"))
// 	}

// 	email, pass, ok := r.BasicAuth()
// 	if !ok {
// 		return auth.NewAuthError("must provide email and password in Basic auth")
// 	}

// 	addr, err := mail.ParseAddress(email)
// 	if err != nil {
// 		return auth.NewAuthError("invalid email format")
// 	}

// 	usr, err := h.user.Authenticate(ctx, *addr, pass)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, user.ErrNotFound):
// 			return v1.NewTrustedError(err, http.StatusNotFound)
// 		case errors.Is(err, user.ErrAuthenticationFailure):
// 			return auth.NewAuthError(err.Error())
// 		default:
// 			return fmt.Errorf("authenticate: %w", err)
// 		}
// 	}

// 	claims := auth.Claims{
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			Subject:   usr.ID.String(),
// 			Issuer:    "service project",
// 			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
// 			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
// 		},
// 		Roles: usr.Roles,
// 	}

// 	token, err := h.auth.GenerateToken(kid, claims)
// 	if err != nil {
// 		return fmt.Errorf("generatetoken: %w", err)
// 	}

// 	return web.Respond(ctx, w, toToken(token), http.StatusOK)
// }
