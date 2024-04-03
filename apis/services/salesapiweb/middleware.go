package salesapiweb

import (
	"context"

	eauth "encore.dev/beta/auth"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/api/mid"
)

// =============================================================================
// JWT or Basic Athentication handling

//lint:ignore U1000 "called by encore"
//encore:authhandler
func (s *Service) AuthHandler(ctx context.Context, ap *mid.AuthParams) (eauth.UID, *auth.Claims, error) {
	return mid.AuthHandler(ctx, s.auth, s.busCrud.user, ap)
}

// =============================================================================
// Global middleware functions
// The order matters so be careful when injecting new middleware.

//lint:ignore U1000 "called by encore"
//encore:middleware target=all
func (s *Service) panics(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.Panics(s.mtrcs, req, next)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:metrics
func (s *Service) metrics(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.Metrics(s.mtrcs, req, next)
}

// =============================================================================
// Authorization related middleware
// These middleware functions must come after the global middleware functions
// above. These are targeted so the order doesn't matter.

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize
func (s *Service) authorize(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.Authorize(s.auth, req, next)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_user
func (s *Service) authorizeUser(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.AuthorizeUser(s.auth, s.busCrud.user, req, next)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_product
func (s *Service) authorizeProduct(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.AuthorizeProduct(s.auth, s.busCrud.product, req, next)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_home
func (s *Service) authorizeHome(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.AuthorizeHome(s.auth, s.busCrud.home, req, next)
}
