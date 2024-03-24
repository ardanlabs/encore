// Package tranapi maintains the group of handlers for transaction example.
package tranapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
)

// Handlers manages the set of handler functions for this domain.
type Handlers struct {
	user    *user.Core
	product *product.Core
}

// New constructs a Handlers for use.
func New(user *user.Core, product *product.Core) *Handlers {
	return &Handlers{
		user:    user,
		product: product,
	}
}

// Create adds a new user and product at the same time under a single transaction.
func (h *Handlers) Create(ctx context.Context, app AppNewTran) (AppProduct, error) {
	h, err := h.executeUnderTransaction(ctx)
	if err != nil {
		return AppProduct{}, errs.New(http.StatusInternalServerError, err)
	}

	np, err := toCoreNewProduct(app.Product)
	if err != nil {
		return AppProduct{}, errs.New(http.StatusBadRequest, err)
	}

	nu, err := toCoreNewUser(app.User)
	if err != nil {
		return AppProduct{}, errs.New(http.StatusBadRequest, err)
	}

	usr, err := h.user.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return AppProduct{}, errs.New(http.StatusConflict, user.ErrUniqueEmail)
		}
		return AppProduct{}, errs.Newf(http.StatusInternalServerError, "create: usr[%+v]: %s", usr, err)
	}

	np.UserID = usr.ID

	prd, err := h.product.Create(ctx, np)
	if err != nil {
		return AppProduct{}, errs.Newf(http.StatusInternalServerError, "create: prd[%+v]: %s", prd, err)
	}

	return toAppProduct(prd), nil
}
