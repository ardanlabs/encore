// Package productapi maintains the group of handlers for product access.
package productapi

import (
	"context"
	"errors"
	"net/http"

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
		return AppProduct{}, errs.New(http.StatusBadRequest, err)
	}

	prd, err := api.product.Create(ctx, np)
	if err != nil {
		return AppProduct{}, errs.Newf(http.StatusInternalServerError, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}

// Update updates an existing product.
func (api *API) Update(ctx context.Context, productID string, app AppUpdateProduct) (AppProduct, error) {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return AppProduct{}, errs.Newf(http.StatusInternalServerError, "product missing in context: %s", err)
	}

	updPrd, err := api.product.Update(ctx, prd, toCoreUpdateProduct(app))
	if err != nil {
		return AppProduct{}, errs.Newf(http.StatusInternalServerError, "update: productID[%s] up[%+v]: %s", prd.ID, app, err)
	}

	return toAppProduct(updPrd), nil
}

// Delete removes a product from the system.
func (api *API) Delete(ctx context.Context, productID string) error {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return errs.Newf(http.StatusInternalServerError, "productID[%s] missing in context: %s", productID, err)
	}

	if err := api.product.Delete(ctx, prd); err != nil {
		return errs.Newf(http.StatusInternalServerError, "delete: productID[%s]: %s", prd.ID, err)
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
		return page.Document[AppProduct]{}, errs.Newf(http.StatusInternalServerError, "query: %s", err)
	}

	total, err := api.product.Count(ctx, filter)
	if err != nil {
		return page.Document[AppProduct]{}, errs.Newf(http.StatusInternalServerError, "count: %s", err)
	}

	return page.NewDocument(toAppProducts(prds), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a product by its ID.
func (api *API) QueryByID(ctx context.Context, productID string) (AppProduct, error) {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return AppProduct{}, errs.Newf(http.StatusInternalServerError, "querybyid: productID[%s]: %s", productID, err)
	}

	return toAppProduct(prd), nil
}
