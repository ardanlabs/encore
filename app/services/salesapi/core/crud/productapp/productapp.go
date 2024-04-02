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

// Core manages the set of handler functions for this domain.
type Core struct {
	product *product.Core
}

// New constructs a Handlers for use.
func New(product *product.Core) *Core {
	return &Core{
		product: product,
	}
}

// Create adds a new product to the system.
func (c *Core) Create(ctx context.Context, app AppNewProduct) (AppProduct, error) {
	np, err := toCoreNewProduct(ctx, app)
	if err != nil {
		return AppProduct{}, errs.New(eerrs.FailedPrecondition, err)
	}

	prd, err := c.product.Create(ctx, np)
	if err != nil {
		return AppProduct{}, errs.Newf(eerrs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}

// Update updates an existing product.
func (c *Core) Update(ctx context.Context, productID string, app AppUpdateProduct) (AppProduct, error) {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return AppProduct{}, errs.Newf(eerrs.Internal, "product missing in context: %s", err)
	}

	updPrd, err := c.product.Update(ctx, prd, toCoreUpdateProduct(app))
	if err != nil {
		return AppProduct{}, errs.Newf(eerrs.Internal, "update: productID[%s] up[%+v]: %s", prd.ID, app, err)
	}

	return toAppProduct(updPrd), nil
}

// Delete removes a product from the system.
func (c *Core) Delete(ctx context.Context, productID string) error {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return errs.Newf(eerrs.Internal, "productID[%s] missing in context: %s", productID, err)
	}

	if err := c.product.Delete(ctx, prd); err != nil {
		return errs.Newf(eerrs.Internal, "delete: productID[%s]: %s", prd.ID, err)
	}

	return nil
}

// Query returns a list of products with paging.
func (c *Core) Query(ctx context.Context, qp QueryParams) (page.Document[AppProduct], error) {
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

	prds, err := c.product.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[AppProduct]{}, errs.Newf(eerrs.Internal, "query: %s", err)
	}

	total, err := c.product.Count(ctx, filter)
	if err != nil {
		return page.Document[AppProduct]{}, errs.Newf(eerrs.Internal, "count: %s", err)
	}

	return page.NewDocument(toAppProducts(prds), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a product by its ID.
func (c *Core) QueryByID(ctx context.Context, productID string) (AppProduct, error) {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return AppProduct{}, errs.Newf(eerrs.Internal, "querybyid: productID[%s]: %s", productID, err)
	}

	return toAppProduct(prd), nil
}
