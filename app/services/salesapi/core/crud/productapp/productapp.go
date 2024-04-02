// Package productapp maintains the app layer api for the product domain.
package productapp

import (
	"context"
	"errors"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/api/mid"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/core/crud/product"
)

// Set of error variables for handling product group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// API manages the set of handler functions for this domain.
type API struct {
	product *product.Core
}

// New constructs a Handlers for use.
func New(product *product.Core) *API {
	return &API{
		product: product,
	}
}

// Create adds a new product to the system.
func (api *API) Create(ctx context.Context, app AppNewProduct) (AppProduct, error) {
	np, err := toCoreNewProduct(ctx, app)
	if err != nil {
		return AppProduct{}, errs.New(eerrs.FailedPrecondition, err)
	}

	prd, err := api.product.Create(ctx, np)
	if err != nil {
		return AppProduct{}, errs.Newf(eerrs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}

// Update updates an existing product.
func (api *API) Update(ctx context.Context, productID string, app AppUpdateProduct) (AppProduct, error) {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return AppProduct{}, errs.Newf(eerrs.Internal, "product missing in context: %s", err)
	}

	updPrd, err := api.product.Update(ctx, prd, toCoreUpdateProduct(app))
	if err != nil {
		return AppProduct{}, errs.Newf(eerrs.Internal, "update: productID[%s] up[%+v]: %s", prd.ID, app, err)
	}

	return toAppProduct(updPrd), nil
}

// Delete removes a product from the system.
func (api *API) Delete(ctx context.Context, productID string) error {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return errs.Newf(eerrs.Internal, "productID[%s] missing in context: %s", productID, err)
	}

	if err := api.product.Delete(ctx, prd); err != nil {
		return errs.Newf(eerrs.Internal, "delete: productID[%s]: %s", prd.ID, err)
	}

	return nil
}

// Query returns a list of products with paging.
func (api *API) Query(ctx context.Context, qp QueryParams) (page.Document[AppProduct], error) {
	if err := validatePaging(qp); err != nil {
		return page.Document[AppProduct]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Document[AppProduct]{}, err
	}

	orderBy, err := parseOrder(qp)
	if err != nil {
		return page.Document[AppProduct]{}, err
	}

	prds, err := api.product.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[AppProduct]{}, errs.Newf(eerrs.Internal, "query: %s", err)
	}

	total, err := api.product.Count(ctx, filter)
	if err != nil {
		return page.Document[AppProduct]{}, errs.Newf(eerrs.Internal, "count: %s", err)
	}

	return page.NewDocument(toAppProducts(prds), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a product by its ID.
func (api *API) QueryByID(ctx context.Context, productID string) (AppProduct, error) {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return AppProduct{}, errs.Newf(eerrs.Internal, "querybyid: productID[%s]: %s", productID, err)
	}

	return toAppProduct(prd), nil
}
