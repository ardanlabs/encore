package encore

import (
	"context"

	encauth "encore.dev/beta/auth"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/web/auth"
	"github.com/ardanlabs/encore/business/web/mid"
)

// =============================================================================
// JWT or Basic Athentication handling

//lint:ignore U1000 "called by encore"
//encore:authhandler
func (s *Service) authHandler(ctx context.Context, ap *mid.AuthParams) (encauth.UID, *auth.Claims, error) {
	return mid.AuthHandler(ctx, s.Auth, s.UsrCore, ap)
}

// =============================================================================
// Global middleware functions
// The order matters so be careful when injecting new middleware.

//lint:ignore U1000 "called by encore"
//encore:middleware target=all
func (s *Service) metrics(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.Metrics(s.Metrics, req, next)
}

// =============================================================================
// Authorization related middleware
// These middleware functions must come after the global middleware functions
// above. These are targeted so the order doesn't matter.

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_admin_only
func (s *Service) authorizeAdminOnly(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.AuthorizeAdminOnly(s.Auth, req, next)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_user
func (s *Service) authorizeUser(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.AuthorizeUser(s.Auth, s.UsrCore, req, next)
}
