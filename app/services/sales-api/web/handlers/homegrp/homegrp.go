// Package homegrp maintains the group of handlers for home access.
package homegrp

import (
	"context"
	"errors"
	"net/http"

	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/business/web/mid"
	"github.com/ardanlabs/encore/business/web/page"
)

// Set of error variables for handling home group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Handlers manages the set of handler functions for this domain.
type Handlers struct {
	home *home.Core
}

// New constructs a Handlers for use.
func New(home *home.Core) *Handlers {
	return &Handlers{
		home: home,
	}
}

// Create adds a new home to the system.
func (h *Handlers) Create(ctx context.Context, app AppNewHome) (AppHome, error) {
	nh, err := toCoreNewHome(ctx, app)
	if err != nil {
		return AppHome{}, errs.New(http.StatusBadRequest, err)
	}

	hme, err := h.home.Create(ctx, nh)
	if err != nil {
		return AppHome{}, errs.Newf(http.StatusInternalServerError, "create: hme[%+v]: %s", app, err)
	}

	return toAppHome(hme), nil
}

// Update updates an existing home.
func (h *Handlers) Update(ctx context.Context, userID string, app AppUpdateHome) (AppHome, error) {
	uh, err := toCoreUpdateHome(app)
	if err != nil {
		return AppHome{}, errs.New(http.StatusBadRequest, err)
	}

	hme, err := mid.GetHome(ctx)
	if err != nil {
		return AppHome{}, errs.Newf(http.StatusInternalServerError, "home missing in context: %s", err)
	}

	updUsr, err := h.home.Update(ctx, hme, uh)
	if err != nil {
		return AppHome{}, errs.Newf(http.StatusInternalServerError, "update: homeID[%s] uh[%+v]: %s", hme.ID, uh, err)
	}

	return toAppHome(updUsr), nil
}

// Delete removes a home from the system.
func (h *Handlers) Delete(ctx context.Context, homeID string) error {
	hme, err := mid.GetHome(ctx)
	if err != nil {
		return errs.Newf(http.StatusInternalServerError, "homeID[%s] missing in context: %s", homeID, err)
	}

	if err := h.home.Delete(ctx, hme); err != nil {
		return errs.Newf(http.StatusInternalServerError, "delete: homeID[%s]: %s", hme.ID, err)
	}

	return nil
}

// Query returns a list of homes with paging.
func (h *Handlers) Query(ctx context.Context, qp QueryParams) (page.Document[AppHome], error) {
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

	hmes, err := h.home.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[AppHome]{}, errs.Newf(http.StatusInternalServerError, "query: %s", err)
	}

	total, err := h.home.Count(ctx, filter)
	if err != nil {
		return page.Document[AppHome]{}, errs.Newf(http.StatusInternalServerError, "count: %s", err)
	}

	return page.NewDocument(toAppHomes(hmes), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a home by its ID.
func (h *Handlers) QueryByID(ctx context.Context, homeID string) (AppHome, error) {
	hme, err := mid.GetHome(ctx)
	if err != nil {
		return AppHome{}, errs.Newf(http.StatusInternalServerError, "querybyid: homeID[%s]: %s", homeID, err)
	}

	return toAppHome(hme), nil
}
