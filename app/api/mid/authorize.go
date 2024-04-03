package mid

import (
	"errors"

	eauth "encore.dev/beta/auth"
	eerrs "encore.dev/beta/errs"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/core/crud/homebus"
	"github.com/ardanlabs/encore/business/core/crud/productbus"
	"github.com/ardanlabs/encore/business/core/crud/userbus"
	"github.com/google/uuid"
)

// Authorize checks the user making the request is an admin or user.
func Authorize(a *auth.Auth, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	claims := eauth.Data().(*auth.Claims)

	rule := auth.RuleAdminOnly
	for _, tag := range req.Data().API.Tags {
		switch tag {
		case "as_any_role":
			rule = auth.RuleAny
		case "as_user_role":
			rule = auth.RuleUserOnly
		}
	}

	if err := a.Authorize(ctx, *claims, uuid.UUID{}, rule); err != nil {
		return errs.NewResponsef(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
	}

	return next(req)
}

// AuthorizeUser checks the user making the call has specified a user id on
// the route that matches the claims.
func AuthorizeUser(a *auth.Auth, userBus *userbus.Core, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	rule := auth.RuleAdminOrSubject
	for _, tag := range req.Data().API.Tags {
		if tag == "as_admin_role" {
			rule = auth.RuleAdminOnly
			break
		}
	}

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		var err error
		userID, err = uuid.Parse(id.Value)
		if err != nil {
			return errs.NewResponse(eerrs.Unauthenticated, ErrInvalidID)
		}

		usr, err := userBus.QueryByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, userbus.ErrNotFound):
				return errs.NewResponse(eerrs.Unauthenticated, err)

			default:
				return errs.NewResponsef(eerrs.Unauthenticated, "querybyid: userID[%s]: %s", userID, err)
			}
		}

		req = setUser(req, usr)
	}

	claims := eauth.Data().(*auth.Claims)

	if err := a.Authorize(ctx, *claims, userID, rule); err != nil {
		return errs.NewResponsef(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, rule, err)
	}

	return next(req)
}

// AuthorizeProduct checks the user making the call has specified a product id on
// the route that matches the claims.
func AuthorizeProduct(a *auth.Auth, productCore *productbus.Core, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		productID, err := uuid.Parse(id.Value)
		if err != nil {
			return errs.NewResponse(eerrs.Unauthenticated, ErrInvalidID)
		}

		prd, err := productCore.QueryByID(ctx, productID)
		if err != nil {
			switch {
			case errors.Is(err, productbus.ErrNotFound):
				return errs.NewResponse(eerrs.Unauthenticated, err)

			default:
				return errs.NewResponsef(eerrs.Internal, "querybyid: productID[%s]: %s", productID, err)
			}
		}

		userID = prd.UserID
		req = setProduct(req, prd)
	}

	claims := eauth.Data().(*auth.Claims)

	if err := a.Authorize(req.Context(), *claims, userID, auth.RuleAdminOrSubject); err != nil {
		return errs.NewResponsef(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err)
	}

	return next(req)
}

// AuthorizeHome checks the user making the call has specified a home id on
// the route that matches the claims.
func AuthorizeHome(a *auth.Auth, homeCore *homebus.Core, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		homeID, err := uuid.Parse(id.Value)
		if err != nil {
			return errs.NewResponse(eerrs.Unauthenticated, ErrInvalidID)
		}

		hme, err := homeCore.QueryByID(ctx, homeID)
		if err != nil {
			switch {
			case errors.Is(err, homebus.ErrNotFound):
				return errs.NewResponse(eerrs.Unauthenticated, err)

			default:
				return errs.NewResponsef(eerrs.Unauthenticated, "querybyid: homeID[%s]: %s", homeID, err)
			}
		}

		userID = hme.UserID
		req = setHome(req, hme)
	}

	claims := eauth.Data().(*auth.Claims)

	if err := a.Authorize(req.Context(), *claims, userID, auth.RuleAdminOrSubject); err != nil {
		return errs.NewResponsef(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err)
	}

	return next(req)
}
