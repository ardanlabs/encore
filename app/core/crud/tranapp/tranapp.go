// Package tranapp maintains the app layer api for the tran domain.
package tranapp

import (
	"context"
	"errors"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

// Core manages the set of handler functions for this domain.
type Core struct {
	user    *user.Core
	product *product.Core
}

// NewCore constructs a tran core API for use.
func NewCore(user *user.Core, product *product.Core) *Core {
	return &Core{
		user:    user,
		product: product,
	}
}

// Create adds a new user and product at the same time under a single transaction.
func (c *Core) Create(ctx context.Context, app NewTran) (Product, error) {
	h, err := c.executeUnderTransaction(ctx)
	if err != nil {
		return Product{}, errs.New(eerrs.Internal, err)
	}

	np, err := toBusNewProduct(app.Product)
	if err != nil {
		return Product{}, errs.New(eerrs.FailedPrecondition, err)
	}

	nu, err := toBusNewUser(app.User)
	if err != nil {
		return Product{}, errs.New(eerrs.FailedPrecondition, err)
	}

	usr, err := h.user.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return Product{}, errs.New(eerrs.Aborted, user.ErrUniqueEmail)
		}
		return Product{}, errs.Newf(eerrs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	np.UserID = usr.ID

	prd, err := h.product.Create(ctx, np)
	if err != nil {
		return Product{}, errs.Newf(eerrs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}
