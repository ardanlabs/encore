// Package homeapp maintains the app layer api for the home domain.
package homeapp

import (
	"context"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/api/mid"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/core/crud/home"
)

// Core manages the set of app layer api functions for the home domain.
type Core struct {
	home *home.Core
}

// NewCore constructs a home core API for use.
func NewCore(home *home.Core) *Core {
	return &Core{
		home: home,
	}
}

// Create adds a new home to the system.
func (c *Core) Create(ctx context.Context, app NewHome) (Home, error) {
	nh, err := toBusNewHome(ctx, app)
	if err != nil {
		return Home{}, errs.New(eerrs.FailedPrecondition, err)
	}

	hme, err := c.home.Create(ctx, nh)
	if err != nil {
		return Home{}, errs.Newf(eerrs.Internal, "create: hme[%+v]: %s", app, err)
	}

	return toAppHome(hme), nil
}

// Update updates an existing home.
func (c *Core) Update(ctx context.Context, userID string, app UpdateHome) (Home, error) {
	uh, err := toBusUpdateHome(app)
	if err != nil {
		return Home{}, errs.New(eerrs.FailedPrecondition, err)
	}

	hme, err := mid.GetHome(ctx)
	if err != nil {
		return Home{}, errs.Newf(eerrs.Internal, "home missing in context: %s", err)
	}

	updUsr, err := c.home.Update(ctx, hme, uh)
	if err != nil {
		return Home{}, errs.Newf(eerrs.Internal, "update: homeID[%s] uh[%+v]: %s", hme.ID, uh, err)
	}

	return toAppHome(updUsr), nil
}

// Delete removes a home from the system.
func (c *Core) Delete(ctx context.Context, homeID string) error {
	hme, err := mid.GetHome(ctx)
	if err != nil {
		return errs.Newf(eerrs.Internal, "homeID[%s] missing in context: %s", homeID, err)
	}

	if err := c.home.Delete(ctx, hme); err != nil {
		return errs.Newf(eerrs.Internal, "delete: homeID[%s]: %s", hme.ID, err)
	}

	return nil
}

// Query returns a list of homes with paging.
func (c *Core) Query(ctx context.Context, qp QueryParams) (page.Document[Home], error) {
	if err := validatePaging(qp); err != nil {
		return page.Document[Home]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Document[Home]{}, err
	}

	orderBy, err := parseOrder(qp)
	if err != nil {
		return page.Document[Home]{}, err
	}

	hmes, err := c.home.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[Home]{}, errs.Newf(eerrs.Internal, "query: %s", err)
	}

	total, err := c.home.Count(ctx, filter)
	if err != nil {
		return page.Document[Home]{}, errs.Newf(eerrs.Internal, "count: %s", err)
	}

	return page.NewDocument(toAppHomes(hmes), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a home by its ID.
func (c *Core) QueryByID(ctx context.Context, homeID string) (Home, error) {
	hme, err := mid.GetHome(ctx)
	if err != nil {
		return Home{}, errs.Newf(eerrs.Internal, "querybyid: homeID[%s]: %s", homeID, err)
	}

	return toAppHome(hme), nil
}
