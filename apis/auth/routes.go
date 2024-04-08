package auth

import (
	"context"

	eauth "encore.dev/beta/auth"
	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/app/core/crud/userapp"
	"github.com/ardanlabs/encore/business/api/auth"
	"github.com/ardanlabs/encore/business/api/mid"
)

// =============================================================================
// JWT or Basic Athentication handling

//lint:ignore U1000 "called by encore"
//encore:authhandler
func (s *Service) AuthHandler(ctx context.Context, ap *mid.AuthParams) (eauth.UID, *auth.Claims, error) {
	return mid.AuthHandler(ctx, s.auth, s.userBus, ap)
}

// =============================================================================
// Auth related APIs

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/token/:kid
func (s *Service) UserToken(ctx context.Context, kid string) (userapp.Token, error) {
	return s.userapp.Token(ctx, kid)
}

//lint:ignore U1000 "called by encore"
//encore:api public method=POST path=/v1/authorize
func (s *Service) Authorize(ctx context.Context, authInfo mid.AuthInfo) error {
	if err := s.auth.Authorize(ctx, authInfo.Claims, authInfo.UserID, authInfo.Rule); err != nil {
		return errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", authInfo.Claims.Roles, authInfo.Rule, err)
	}

	return nil
}
