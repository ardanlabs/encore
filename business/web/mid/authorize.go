package mid

import (
	"errors"
	"net/http"

	encauth "encore.dev/beta/auth"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/web/auth"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/google/uuid"
)

// AuthorizeAdminOnly checks the user making the request is an admin.
func AuthorizeAdminOnly(a *auth.Auth, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	claims := encauth.Data().(*auth.Claims)

	if err := a.Authorize(ctx, *claims, uuid.UUID{}, auth.RuleAdminOnly); err != nil {
		return errs.NewResponsef(http.StatusBadRequest, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOnly, err)
	}

	return next(req)
}

// AuthorizeUser checks the user making the call has specified a user id on
// the route that matches the claims.
func AuthorizeUser(a *auth.Auth, usrCore *user.Core, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		userID, err := uuid.Parse(id.Value)
		if err != nil {
			return errs.NewResponse(http.StatusBadRequest, ErrInvalidID)
		}

		usr, err := usrCore.QueryByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrNotFound):
				return errs.NewResponse(http.StatusBadRequest, err)

			default:
				return errs.NewResponsef(http.StatusInternalServerError, "querybyid: userID[%s]: %s", userID, err)
			}
		}

		req = setUser(req, usr)
	}

	claims := encauth.Data().(*auth.Claims)

	if err := a.Authorize(ctx, *claims, userID, auth.RuleAdminOrSubject); err != nil {
		return errs.NewResponsef(http.StatusBadRequest, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err)
	}

	return next(req)
}

// AuthorizeProduct checks the user making the call has specified a user id on
// the route that matches the claims.
func AuthorizeProduct(a *auth.Auth, prdCore *product.Core, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		productID, err := uuid.Parse(id.Value)
		if err != nil {
			return errs.NewResponse(http.StatusBadRequest, ErrInvalidID)
		}

		prd, err := prdCore.QueryByID(ctx, productID)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrNotFound):
				return errs.NewResponse(http.StatusBadRequest, err)

			default:
				return errs.NewResponsef(http.StatusInternalServerError, "querybyid: productID[%s]: %s", productID, err)
			}
		}

		userID = prd.UserID
		ctx = setProduct(ctx, prd)
	}

	claims := encauth.Data().(*auth.Claims)

	if err := a.Authorize(ctx, *claims, userID, auth.RuleAdminOrSubject); err != nil {
		return errs.NewResponsef(http.StatusBadRequest, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err)
	}

	return next(req)
}
