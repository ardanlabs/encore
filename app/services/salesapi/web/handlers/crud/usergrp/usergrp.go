// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"net/http"

	eauth "encore.dev/beta/auth"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/web/auth"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/business/web/mid"
	"github.com/ardanlabs/encore/business/web/page"
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
	claims := eauth.Data().(*auth.Claims)

	tkn, err := h.auth.GenerateToken(kid, *claims)
	if err != nil {
		return Token{}, errs.New(http.StatusInternalServerError, err)
	}

	return toToken(tkn), nil
}

// Create adds a new user to the system.
func (h *Handlers) Create(ctx context.Context, app AppNewUser) (AppUser, error) {
	nc, err := toCoreNewUser(app)
	if err != nil {
		return AppUser{}, errs.New(http.StatusBadRequest, err)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return AppUser{}, errs.New(http.StatusConflict, user.ErrUniqueEmail)
		}
		return AppUser{}, errs.Newf(http.StatusInternalServerError, "create: usr[%+v]: %s", usr, err)
	}

	return toAppUser(usr), nil
}

// Update updates an existing user.
func (h *Handlers) Update(ctx context.Context, userID string, app AppUpdateUser) (AppUser, error) {
	uu, err := toCoreUpdateUser(app)
	if err != nil {
		return AppUser{}, errs.New(http.StatusBadRequest, err)
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		return AppUser{}, errs.Newf(http.StatusInternalServerError, "user missing in context: %s", err)
	}

	updUsr, err := h.user.Update(ctx, usr, uu)
	if err != nil {
		return AppUser{}, errs.Newf(http.StatusInternalServerError, "update: userID[%s] uu[%+v]: %s", usr.ID, uu, err)
	}

	return toAppUser(updUsr), nil
}

// Delete removes a user from the system.
func (h *Handlers) Delete(ctx context.Context, userID string) error {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return errs.Newf(http.StatusInternalServerError, "userID[%s] missing in context: %s", userID, err)
	}

	if err := h.user.Delete(ctx, usr); err != nil {
		return errs.Newf(http.StatusInternalServerError, "delete: userID[%s]: %s", usr.ID, err)
	}

	return nil
}

// Query returns a list of users with paging.
func (h *Handlers) Query(ctx context.Context, qp QueryParams) (page.Document[AppUser], error) {
	if err := validatePaging(qp); err != nil {
		return page.Document[AppUser]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Document[AppUser]{}, err
	}

	orderBy, err := parseOrder(qp)
	if err != nil {
		return page.Document[AppUser]{}, err
	}

	usrs, err := h.user.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[AppUser]{}, errs.Newf(http.StatusInternalServerError, "query: %s", err)
	}

	total, err := h.user.Count(ctx, filter)
	if err != nil {
		return page.Document[AppUser]{}, errs.Newf(http.StatusInternalServerError, "count: %s", err)
	}

	return page.NewDocument(toAppUsers(usrs), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a user by its ID.
func (h *Handlers) QueryByID(ctx context.Context, userID string) (AppUser, error) {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return AppUser{}, errs.Newf(http.StatusInternalServerError, "querybyid: userID[%s]: %s", userID, err)
	}

	return toAppUser(usr), nil
}
