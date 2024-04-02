// Package userapp maintains the app layer api for the user domain.
package userapp

import (
	"context"
	"errors"

	eauth "encore.dev/beta/auth"
	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/api/mid"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

// Core manages the set of handler functions for this domain.
type Core struct {
	user *user.Core
	auth *auth.Auth
}

// New constructs a Handlers for use.
func New(user *user.Core, auth *auth.Auth) *Core {
	return &Core{
		user: user,
		auth: auth,
	}
}

// Token provides an API token for the authenticated user.
func (c *Core) Token(ctx context.Context, kid string) (Token, error) {
	claims := eauth.Data().(*auth.Claims)

	tkn, err := c.auth.GenerateToken(kid, *claims)
	if err != nil {
		return Token{}, errs.New(eerrs.Internal, err)
	}

	return toToken(tkn), nil
}

// Create adds a new user to the system.
func (c *Core) Create(ctx context.Context, app NewUser) (User, error) {
	nc, err := toBusNewUser(app)
	if err != nil {
		return User{}, errs.New(eerrs.FailedPrecondition, err)
	}

	usr, err := c.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return User{}, errs.New(eerrs.Aborted, user.ErrUniqueEmail)
		}
		return User{}, errs.Newf(eerrs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	return toAppUser(usr), nil
}

// Update updates an existing user.
func (c *Core) Update(ctx context.Context, userID string, app UpdateUser) (User, error) {
	uu, err := toBusUpdateUser(app)
	if err != nil {
		return User{}, errs.New(eerrs.FailedPrecondition, err)
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		return User{}, errs.Newf(eerrs.Internal, "user missing in context: %s", err)
	}

	updUsr, err := c.user.Update(ctx, usr, uu)
	if err != nil {
		return User{}, errs.Newf(eerrs.Internal, "update: userID[%s] uu[%+v]: %s", usr.ID, uu, err)
	}

	return toAppUser(updUsr), nil
}

// UpdateRole updates an existing user's role.
func (c *Core) UpdateRole(ctx context.Context, userID string, app UpdateUserRole) (User, error) {
	uu, err := toBusUpdateUserRole(app)
	if err != nil {
		return User{}, errs.New(eerrs.FailedPrecondition, err)
	}

	usr, err := mid.GetUser(ctx)
	if err != nil {
		return User{}, errs.Newf(eerrs.Internal, "user missing in context: %s", err)
	}

	updUsr, err := c.user.Update(ctx, usr, uu)
	if err != nil {
		return User{}, errs.Newf(eerrs.Internal, "updaterole: userID[%s] uu[%+v]: %s", usr.ID, uu, err)
	}

	return toAppUser(updUsr), nil
}

// Delete removes a user from the system.
func (c *Core) Delete(ctx context.Context, userID string) error {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return errs.Newf(eerrs.Internal, "userID[%s] missing in context: %s", userID, err)
	}

	if err := c.user.Delete(ctx, usr); err != nil {
		return errs.Newf(eerrs.Internal, "delete: userID[%s]: %s", usr.ID, err)
	}

	return nil
}

// Query returns a list of users with paging.
func (c *Core) Query(ctx context.Context, qp QueryParams) (page.Document[User], error) {
	if err := validatePaging(qp); err != nil {
		return page.Document[User]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Document[User]{}, err
	}

	orderBy, err := parseOrder(qp)
	if err != nil {
		return page.Document[User]{}, err
	}

	usrs, err := c.user.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[User]{}, errs.Newf(eerrs.Internal, "query: %s", err)
	}

	total, err := c.user.Count(ctx, filter)
	if err != nil {
		return page.Document[User]{}, errs.Newf(eerrs.Internal, "count: %s", err)
	}

	return page.NewDocument(toAppUsers(usrs), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a user by its ID.
func (c *Core) QueryByID(ctx context.Context, userID string) (User, error) {
	usr, err := mid.GetUser(ctx)
	if err != nil {
		return User{}, errs.Newf(eerrs.Internal, "querybyid: userID[%s]: %s", userID, err)
	}

	return toAppUser(usr), nil
}
