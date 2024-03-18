// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	encauth "encore.dev/beta/auth"
	"github.com/ardanlabs/encore/business/core/crud/user"
	v1 "github.com/ardanlabs/encore/business/web/v1"
	"github.com/ardanlabs/encore/business/web/v1/auth"
)

// Handlers manages the set of handler functions for this domain.
type Handlers struct {
	user *user.Core
	auth *auth.Auth
}

// New constructs a Handlers for use.
func New(user *user.Core, auth *auth.Auth) *Handlers {
	return &Handlers{
		user: user,
		auth: auth,
	}
}

// Token provides an API token for the authenticated user.
func (h *Handlers) Token(ctx context.Context, kid string) (Token, error) {
	claims := encauth.Data().(*auth.Claims)

	tkn, err := h.auth.GenerateToken(kid, *claims)
	if err != nil {
		return Token{}, fmt.Errorf("generatetoken: %w", err)
	}

	return toToken(tkn), nil
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

// Update updates an existing user.
func (h *Handlers) Update(ctx context.Context, userID string, app AppUpdateUser) (AppUser, error) {
	uu, err := toCoreUpdateUser(app)
	if err != nil {
		return AppUser{}, v1.NewTrustedError(err, http.StatusBadRequest)
	}

	usr, err := getUser(ctx)
	if err != nil {
		return AppUser{}, fmt.Errorf("user missing in context: %w", err)
	}

	updUsr, err := h.user.Update(ctx, usr, uu)
	if err != nil {
		return AppUser{}, fmt.Errorf("update: userID[%s] uu[%+v]: %w", usr.ID, uu, err)
	}

	return toAppUser(updUsr), nil
}

// Delete removes an existing user.
func (h *Handlers) Delete(ctx context.Context, userID string) error {
	usr, err := getUser(ctx)
	if err != nil {
		return fmt.Errorf("user missing in context: %w", err)
	}

	if err := h.user.Delete(ctx, usr); err != nil {
		return fmt.Errorf("delete: userID[%s]: %w", usr.ID, err)
	}

	return nil
}
