// Package productgrp maintains the group of handlers for product access.
package productgrp

import (
	"context"
	"errors"
	"net/http"

	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/ardanlabs/encore/business/web/mid"
	"github.com/ardanlabs/encore/business/web/page"
)

// Set of error variables for handling product group errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Handlers manages the set of handler functions for this domain.
type Handlers struct {
	product *product.Core
}

// New constructs a Handlers for use.
func New(product *product.Core) *Handlers {
	return &Handlers{
		product: product,
	}
}

// Create adds a new product to the system.
func (h *Handlers) Create(ctx context.Context, app AppNewProduct) (AppProduct, error) {
	np, err := toCoreNewProduct(ctx, app)
	if err != nil {
		return AppProduct{}, errs.New(http.StatusBadRequest, err)
	}

	prd, err := h.product.Create(ctx, np)
	if err != nil {
		return AppProduct{}, errs.Newf(http.StatusInternalServerError, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}

// Update updates an existing product.
func (h *Handlers) Update(ctx context.Context, productID string, app AppUpdateProduct) (AppProduct, error) {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return AppProduct{}, errs.Newf(http.StatusInternalServerError, "product missing in context: %s", err)
	}

	updPrd, err := h.product.Update(ctx, prd, toCoreUpdateProduct(app))
	if err != nil {
		return AppProduct{}, errs.Newf(http.StatusInternalServerError, "update: productID[%s] up[%+v]: %s", prd.ID, app, err)
	}

	return toAppProduct(updPrd), nil
}

// Delete removes a product from the system.
func (h *Handlers) Delete(ctx context.Context, productID string) error {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return errs.Newf(http.StatusInternalServerError, "productID[%s] missing in context: %s", productID, err)
	}

	if err := h.product.Delete(ctx, prd); err != nil {
		return errs.Newf(http.StatusInternalServerError, "delete: productID[%s]: %s", prd.ID, err)
	}

	return nil
}

// Query returns a list of products with paging.
func (h *Handlers) Query(ctx context.Context, qp QueryParams) (page.Document[AppProduct], error) {
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

	prds, err := h.product.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[AppProduct]{}, errs.Newf(http.StatusInternalServerError, "query: %s", err)
	}

	total, err := h.product.Count(ctx, filter)
	if err != nil {
		return page.Document[AppProduct]{}, errs.Newf(http.StatusInternalServerError, "count: %s", err)
	}

	return page.NewDocument(toAppProducts(prds), total, qp.Page, qp.Rows), nil
}

// QueryByID returns a product by its ID.
func (h *Handlers) QueryByID(ctx context.Context, productID string) (AppProduct, error) {
	prd, err := mid.GetProduct(ctx)
	if err != nil {
		return AppProduct{}, errs.Newf(http.StatusInternalServerError, "querybyid: productID[%s]: %s", productID, err)
	}

	return toAppProduct(prd), nil
}
