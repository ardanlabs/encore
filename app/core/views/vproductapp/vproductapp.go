// Package vproductapp maintains the app layer api for the vproduct domain.
package vproductapp

import (
	"context"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/api/page"
	"github.com/ardanlabs/encore/business/core/views/vproduct"
)

// Handlers manages the set of handler functions for this domain.
type Handlers struct {
	vproduct *vproduct.Core
}

// New constructs a Handlers for use.
func New(vproduct *vproduct.Core) *Handlers {
	return &Handlers{
		vproduct: vproduct,
	}
}

// Query returns a list of products with paging.
func (h *Handlers) Query(ctx context.Context, qp QueryParams) (page.Document[Product], error) {
	if err := validatePaging(qp); err != nil {
		return page.Document[Product]{}, err
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return page.Document[Product]{}, err
	}

	orderBy, err := parseOrder(qp)
	if err != nil {
		return page.Document[Product]{}, err
	}

	prds, err := h.vproduct.Query(ctx, filter, orderBy, qp.Page, qp.Rows)
	if err != nil {
		return page.Document[Product]{}, errs.Newf(eerrs.Internal, "query: %s", err)
	}

	total, err := h.vproduct.Count(ctx, filter)
	if err != nil {
		return page.Document[Product]{}, errs.Newf(eerrs.Internal, "count: %s", err)
	}

	return page.NewDocument(toAppProducts(prds), total, qp.Page, qp.Rows), nil
}
