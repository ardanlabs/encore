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
func (h *Handlers) Token(ctx context.Context, kid string) (Token, error) {
	claims := encauth.Data().(*auth.Claims)

	tkn, err := h.auth.GenerateToken(kid, *claims)
	if err != nil {
		return Token{}, fmt.Errorf("generatetoken: %w", err)
	}

	return toToken(tkn), nil
}
