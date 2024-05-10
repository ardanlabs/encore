// Package homeapp maintains the app layer api for the home domain.
package homeapp

import (
	"context"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/sdk/errs"
	"github.com/ardanlabs/encore/app/sdk/mid"
	"github.com/ardanlabs/encore/app/sdk/page"
	"github.com/ardanlabs/encore/business/domain/homebus"
	"github.com/ardanlabs/encore/business/sdk/order"
)

// App manages the set of app layer api functions for the home domain.
type App struct {
	homeBus *homebus.Business
}

// NewApp constructs a home app API for use.
func NewApp(homeBus *homebus.Business) *App {
	return &App{
		homeBus: homeBus,
	}
}

// Create adds a new home to the system.
func (a *App) Create(ctx context.Context, app NewHome) (Home, error) {
	nh, err := toBusNewHome(ctx, app)
	if err != nil {
		return Home{}, errs.New(eerrs.FailedPrecondition, err)
	}

	hme, err := a.homeBus.Create(ctx, nh)
	if err != nil {
		return Home{}, errs.Newf(eerrs.Internal, "create: hme[%+v]: %s", app, err)
	}

	return toAppHome(hme), nil
}

// Update updates an existing home.
func (a *App) Update(ctx context.Context, userID string, app UpdateHome) (Home, error) {
	uh, err := toBusUpdateHome(app)
	if err != nil {
		return Home{}, errs.New(eerrs.FailedPrecondition, err)
	}

	hme, err := mid.GetHome(ctx)
	if err != nil {
		return Home{}, errs.Newf(eerrs.Internal, "home missing in context: %s", err)
	}

	updUsr, err := a.homeBus.Update(ctx, hme, uh)
	if err != nil {
		return Home{}, errs.Newf(eerrs.Internal, "update: homeID[%s] uh[%+v]: %s", hme.ID, uh, err)
	}

	return toAppHome(updUsr), nil
}

// Delete removes a home from the system.
func (a *App) Delete(ctx context.Context, homeID string) error {
	hme, err := mid.GetHome(ctx)
	if err != nil {
		return errs.Newf(eerrs.Internal, "homeID[%s] missing in context: %s", homeID, err)
	}

	if err := a.homeBus.Delete(ctx, hme); err != nil {
		return errs.Newf(eerrs.Internal, "delete: homeID[%s]: %s", hme.ID, err)
	}

	return nil
}

// Query returns a list of homes with paging.
func (a *App) Query(ctx context.Context, qp QueryParams) (page.Document[Home], error) {
	pg, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return page.Document[Home]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Document[Home]{}, err
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return page.Document[Home]{}, err
	}

	hmes, err := a.homeBus.Query(ctx, filter, orderBy, pg.Number, pg.RowsPerPage)
	if err != nil {
		return page.Document[Home]{}, errs.Newf(eerrs.Internal, "query: %s", err)
	}

	total, err := a.homeBus.Count(ctx, filter)
	if err != nil {
		return page.Document[Home]{}, errs.Newf(eerrs.Internal, "count: %s", err)
	}

	return page.NewDocument(toAppHomes(hmes), total, pg.Number, pg.RowsPerPage), nil
}

// QueryByID returns a home by its ID.
func (a *App) QueryByID(ctx context.Context, homeID string) (Home, error) {
	hme, err := mid.GetHome(ctx)
	if err != nil {
		return Home{}, errs.Newf(eerrs.Internal, "querybyid: homeID[%s]: %s", homeID, err)
	}

	return toAppHome(hme), nil
}
