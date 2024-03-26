package mid

import (
	"errors"
	"net/http"

	eauth "encore.dev/beta/auth"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/api/errs"
	"github.com/ardanlabs/encore/business/core/crud/home"
	"github.com/ardanlabs/encore/business/core/crud/product"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/google/uuid"
)

// AuthorizeAny checks the user making the request is an admin or user.
func AuthorizeAny(a *auth.Auth, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	claims := eauth.Data().(*auth.Claims)

	if err := a.Authorize(ctx, *claims, uuid.UUID{}, auth.RuleAny); err != nil {
		return errs.NewResponsef(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAny, err)
	}

	return next(req)
}

// AuthorizeUserOnly checks the user making the request is a user.
func AuthorizeUserOnly(a *auth.Auth, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	claims := eauth.Data().(*auth.Claims)

	if err := a.Authorize(ctx, *claims, uuid.UUID{}, auth.RuleUserOnly); err != nil {
		return errs.NewResponsef(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleUserOnly, err)
	}

	return next(req)
}

// AuthorizeAdminOnly checks the user making the request is an admin.
func AuthorizeAdminOnly(a *auth.Auth, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	claims := eauth.Data().(*auth.Claims)

	if err := a.Authorize(ctx, *claims, uuid.UUID{}, auth.RuleAdminOnly); err != nil {
		return errs.NewResponsef(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOnly, err)
	}

	return next(req)
}

// AuthorizeUser checks the user making the call has specified a user id on
// the route that matches the claims.
func AuthorizeUser(a *auth.Auth, userCore *user.Core, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		var err error
		userID, err = uuid.Parse(id.Value)
		if err != nil {
			return errs.NewResponse(http.StatusUnauthorized, ErrInvalidID)
		}

		usr, err := userCore.QueryByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrNotFound):
				return errs.NewResponse(http.StatusUnauthorized, err)

			default:
				return errs.NewResponsef(http.StatusUnauthorized, "querybyid: userID[%s]: %s", userID, err)
			}
		}

		req = setUser(req, usr)
	}

	claims := eauth.Data().(*auth.Claims)

	if err := a.Authorize(ctx, *claims, userID, auth.RuleAdminOrSubject); err != nil {
		return errs.NewResponsef(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err)
	}

	return next(req)
}

// AuthorizeProduct checks the user making the call has specified a product id on
// the route that matches the claims.
func AuthorizeProduct(a *auth.Auth, productCore *product.Core, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		productID, err := uuid.Parse(id.Value)
		if err != nil {
			return errs.NewResponse(http.StatusUnauthorized, ErrInvalidID)
		}

		prd, err := productCore.QueryByID(ctx, productID)
		if err != nil {
			switch {
			case errors.Is(err, product.ErrNotFound):
				return errs.NewResponse(http.StatusUnauthorized, err)

			default:
				return errs.NewResponsef(http.StatusInternalServerError, "querybyid: productID[%s]: %s", productID, err)
			}
		}

		userID = prd.UserID
		req = setProduct(req, prd)
	}

	claims := eauth.Data().(*auth.Claims)

	if err := a.Authorize(req.Context(), *claims, userID, auth.RuleAdminOrSubject); err != nil {
		return errs.NewResponsef(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err)
	}

	return next(req)
}

// AuthorizeHome checks the user making the call has specified a home id on
// the route that matches the claims.
func AuthorizeHome(a *auth.Auth, homeCore *home.Core, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		homeID, err := uuid.Parse(id.Value)
		if err != nil {
			return errs.NewResponse(http.StatusUnauthorized, ErrInvalidID)
		}

		hme, err := homeCore.QueryByID(ctx, homeID)
		if err != nil {
			switch {
			case errors.Is(err, home.ErrNotFound):
				return errs.NewResponse(http.StatusUnauthorized, err)

			default:
				return errs.NewResponsef(http.StatusUnauthorized, "querybyid: homeID[%s]: %s", homeID, err)
			}
		}

		userID = hme.UserID
		req = setHome(req, hme)
	}

	claims := eauth.Data().(*auth.Claims)

	if err := a.Authorize(req.Context(), *claims, userID, auth.RuleAdminOrSubject); err != nil {
		return errs.NewResponsef(http.StatusUnauthorized, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err)
	}

	return next(req)
}
