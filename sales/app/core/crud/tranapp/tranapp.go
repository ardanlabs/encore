// Package tranapp maintains the app layer api for the tran domain.
package tranapp

import (
	"context"
	"errors"

	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/sales/app/api/errs"
	"github.com/ardanlabs/encore/sales/business/core/crud/productbus"
	"github.com/ardanlabs/encore/sales/business/core/crud/userbus"
)

// Core manages the set of app layer api functions for the tran domain.
type Core struct {
	userBus    *userbus.Core
	productBus *productbus.Core
}

// NewCore constructs a tran core API for use.
func NewCore(userBus *userbus.Core, productBus *productbus.Core) *Core {
	return &Core{
		userBus:    userBus,
		productBus: productBus,
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

	usr, err := h.userBus.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return Product{}, errs.New(eerrs.Aborted, userbus.ErrUniqueEmail)
		}
		return Product{}, errs.Newf(eerrs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	np.UserID = usr.ID

	prd, err := h.productBus.Create(ctx, np)
	if err != nil {
		return Product{}, errs.Newf(eerrs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}
