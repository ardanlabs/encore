// Package homeapi maintains the group of handlers for home access.
package homeapi

import (
	"context"
	"errors"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/api/mid"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/core/crud/home"
)

// Set of error variables for handling home group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// API manages the set of handler functions for this domain.
type API struct {
	home *home.Core
}

// New constructs a Handlers for use.
func New(home *home.Core) *API {
	return &API{
		home: home,
	}
}

// Create adds a new home to the system.
func (api *API) Create(ctx context.Context, app AppNewHome) (AppHome, error) {
	nh, err := toCoreNewHome(ctx, app)
	if err != nil {
		return AppHome{}, errs.New(eerrs.FailedPrecondition, err)
	}

	hme, err := api.home.Create(ctx, nh)
	if err != nil {
		return AppHome{}, errs.Newf(eerrs.Internal, "create: hme[%+v]: %s", app, err)
	}

	return toAppHome(hme), nil
}

// Update updates an existing home.
func (api *API) Update(ctx context.Context, userID string, app AppUpdateHome) (AppHome, error) {
	uh, err := toCoreUpdateHome(app)
	if err != nil {
		return AppHome{}, errs.New(eerrs.FailedPrecondition, err)
	}

	hme, err := mid.GetHome(ctx)
	if err != nil {
		return AppHome{}, errs.Newf(eerrs.Internal, "home missing in context: %s", err)
	}

	updUsr, err := api.home.Update(ctx, hme, uh)
	if err != nil {
		return AppHome{}, errs.Newf(eerrs.Internal, "update: homeID[%s] uh[%+v]: %s", hme.ID, uh, err)
	}

	return toAppHome(updUsr), nil
}

// Delete removes a home from the system.
func (api *API) Delete(ctx context.Context, homeID string) error {
	hme, err := mid.GetHome(ctx)
	if err != nil {
		return errs.Newf(eerrs.Internal, "homeID[%s] missing in context: %s", homeID, err)
	}

	if err := api.home.Delete(ctx, hme); err != nil {
		return errs.Newf(eerrs.Internal, "delete: homeID[%s]: %s", hme.ID, err)
	}

	return nil
}

// Query returns a list of homes with paging.
func (api *API) Query(ctx context.Context, qp QueryParams) (page.Document[AppHome], error) {
	if err := validatePaging(qp); err != nil {
		return page.Document[AppHome]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Document[AppHome]{}, err
	}

	orderBy, err := parseOrder(qp)
	if err != nil {
		return page.Document[AppHome]{}, err
	}

	hmes, err := api.home.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[AppHome]{}, errs.Newf(eerrs.Internal, "query: %s", err)
	}

	total, err := api.home.Count(ctx, filter)
	if err != nil {
		return page.Document[AppHome]{}, errs.Newf(eerrs.Internal, "count: %s", err)
	}

	return page.NewDocument(toAppHomes(hmes), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a home by its ID.
func (api *API) QueryByID(ctx context.Context, homeID string) (AppHome, error) {
	hme, err := mid.GetHome(ctx)
	if err != nil {
		return AppHome{}, errs.Newf(eerrs.Internal, "querybyid: homeID[%s]: %s", homeID, err)
	}

	return toAppHome(hme), nil
}
