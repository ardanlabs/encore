package mid

import (
	"errors"
	"fmt"

	eauth "encore.dev/beta/auth"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/core/crud/homebus"
	"github.com/ardanlabs/encore/business/core/crud/productbus"
	"github.com/ardanlabs/encore/business/core/crud/userbus"
	"github.com/google/uuid"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// Authorize checks the user making the request is an admin or user.
func Authorize(req middleware.Request, next middleware.Next) (AuthInfo, error) {
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

	p := AuthInfo{
		Claims: *claims,
		UserID: uuid.UUID{},
		Rule:   rule,
	}

	return p, nil
}

// AuthorizeUser checks the user making the call has specified a user id on
// the route that matches the claims.
func AuthorizeUser(userBus *userbus.Core, req middleware.Request, next middleware.Next) (AuthInfo, error) {
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
			return AuthInfo{}, ErrInvalidID
		}

		usr, err := userBus.QueryByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, userbus.ErrNotFound):
				return AuthInfo{}, err

			default:
				return AuthInfo{}, fmt.Errorf("querybyid: userID[%s]: %s", userID, err)
			}
		}

		req = setUser(req, usr)
	}

	claims := eauth.Data().(*auth.Claims)

	p := AuthInfo{
		Claims: *claims,
		UserID: userID,
		Rule:   rule,
	}

	return p, nil
}

// AuthorizeProduct checks the user making the call has specified a product id on
// the route that matches the claims.
func AuthorizeProduct(productBus *productbus.Core, req middleware.Request, next middleware.Next) (AuthInfo, error) {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		productID, err := uuid.Parse(id.Value)
		if err != nil {
			return AuthInfo{}, ErrInvalidID
		}

		prd, err := productBus.QueryByID(ctx, productID)
		if err != nil {
			switch {
			case errors.Is(err, productbus.ErrNotFound):
				return AuthInfo{}, err

			default:
				return AuthInfo{}, fmt.Errorf("querybyid: productID[%s]: %s", productID, err)
			}
		}

		userID = prd.UserID
		req = setProduct(req, prd)
	}

	claims := eauth.Data().(*auth.Claims)

	p := AuthInfo{
		Claims: *claims,
		UserID: userID,
		Rule:   auth.RuleAdminOrSubject,
	}

	return p, nil
}

// AuthorizeHome checks the user making the call has specified a home id on
// the route that matches the claims.
func AuthorizeHome(homeBus *homebus.Core, req middleware.Request, next middleware.Next) (AuthInfo, error) {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		homeID, err := uuid.Parse(id.Value)
		if err != nil {
			return AuthInfo{}, ErrInvalidID
		}

		hme, err := homeBus.QueryByID(ctx, homeID)
		if err != nil {
			switch {
			case errors.Is(err, homebus.ErrNotFound):
				return AuthInfo{}, err

			default:
				return AuthInfo{}, fmt.Errorf("querybyid: homeID[%s]: %s", homeID, err)
			}
		}

		userID = hme.UserID
		req = setHome(req, hme)
	}

	claims := eauth.Data().(*auth.Claims)

	p := AuthInfo{
		Claims: *claims,
		UserID: userID,
		Rule:   auth.RuleAdminOrSubject,
	}

	return p, nil
}
