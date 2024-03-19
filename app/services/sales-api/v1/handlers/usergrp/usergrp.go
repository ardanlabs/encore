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
	"github.com/ardanlabs/encore/business/web/v1/mid"
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

	usr, err := mid.GetUser(ctx)
	if err != nil {
		return AppUser{}, fmt.Errorf("user missing in context: %w", err)
	}

	updUsr, err := h.user.Update(ctx, usr, uu)
	if err != nil {
		return AppUser{}, fmt.Errorf("update: userID[%s] uu[%+v]: %w", usr.ID, uu, err)
	}

	return toAppUser(updUsr), nil
}

// Delete removes a user from the system.
func (h *Handlers) Delete(ctx context.Context, userID string) error {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return fmt.Errorf("user missing in context: %w", err)
	}

	if err := h.user.Delete(ctx, usr); err != nil {
		return fmt.Errorf("delete: userID[%s]: %w", usr.ID, err)
	}

	return nil
}

// Query returns a list of users with paging.
func (h *Handlers) Query(ctx context.Context, qp QueryParams) (v1.PageDocument[AppUser], error) {
	if err := validatePaging(qp); err != nil {
		return v1.PageDocument[AppUser]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return v1.PageDocument[AppUser]{}, err
	}

	orderBy, err := parseOrder(qp)
	if err != nil {
		return v1.PageDocument[AppUser]{}, err
	}

	users, err := h.user.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return v1.PageDocument[AppUser]{}, fmt.Errorf("query: %w", err)
	}

	total, err := h.user.Count(ctx, filter)
	if err != nil {
		return v1.PageDocument[AppUser]{}, fmt.Errorf("count: %w", err)
	}

	return v1.NewPageDocument(toAppUsers(users), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a user by its ID.
func (h *Handlers) QueryByID(ctx context.Context) (AppUser, error) {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return AppUser{}, fmt.Errorf("querybyid: %w", err)
	}

	return toAppUser(usr), nil
}
