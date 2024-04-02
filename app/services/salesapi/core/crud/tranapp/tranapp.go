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

// API manages the set of handler functions for this domain.
type API struct {
	user    *user.Core
	product *product.Core
}

// New constructs a Handlers for use.
func New(user *user.Core, product *product.Core) *API {
	return &API{
		user:    user,
		product: product,
	}
}

// Create adds a new user and product at the same time under a single transaction.
func (api *API) Create(ctx context.Context, app AppNewTran) (AppProduct, error) {
	h, err := api.executeUnderTransaction(ctx)
	if err != nil {
		return AppProduct{}, errs.New(eerrs.Internal, err)
	}

	np, err := toCoreNewProduct(app.Product)
	if err != nil {
		return AppProduct{}, errs.New(eerrs.FailedPrecondition, err)
	}

	nu, err := toCoreNewUser(app.User)
	if err != nil {
		return AppProduct{}, errs.New(eerrs.FailedPrecondition, err)
	}

	usr, err := h.user.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return AppProduct{}, errs.New(eerrs.Aborted, user.ErrUniqueEmail)
		}
		return AppProduct{}, errs.Newf(eerrs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	np.UserID = usr.ID

	prd, err := h.product.Create(ctx, np)
	if err != nil {
		return AppProduct{}, errs.Newf(eerrs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}
