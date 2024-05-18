package auth

import (
	"context"
	"strings"

	eauth "encore.dev/beta/auth"
	"github.com/ardanlabs/encore/app/sdk/auth"
	"github.com/ardanlabs/encore/app/sdk/errs"
	"github.com/ardanlabs/encore/app/sdk/mid"
)

// =============================================================================
// JWT or Basic Athentication handling

type authParams struct {
	Authorization string `header:"Authorization"`
}

//lint:ignore U1000 "called by encore"
//encore:authhandler
func (s *Service) AuthHandler(ctx context.Context, ap *authParams) (eauth.UID, *auth.Claims, error) {
	parts := strings.Split(ap.Authorization, " ")

	switch parts[0] {
	case "Bearer":
		return mid.Bearer(ctx, s.auth, ap.Authorization)

	case "Basic":
		return mid.Basic(ctx, s.auth, s.userBus, ap.Authorization)
	}

	return "", nil, errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action")
}

// =============================================================================
// Auth related APIs

type token struct {
	Token string `json:"token"`
}

//lint:ignore U1000 "called by encore"
//encore:api auth method=GET path=/v1/token/:kid
func (s *Service) UserToken(ctx context.Context, kid string) (token, error) {

	// The BearerBasic middleware function generates the claims.
	claims := eauth.Data().(*auth.Claims)

	tkn, err := s.auth.GenerateToken(kid, *claims)
	if err != nil {
		return token{}, errs.New(errs.Internal, err)
	}

	return token{tkn}, nil
}

//lint:ignore U1000 "called by encore"
//encore:api private method=POST path=/v1/authorize
func (s *Service) Authorize(ctx context.Context, authInfo mid.AuthInfo) error {
	if err := s.auth.Authorize(ctx, authInfo.Claims, authInfo.UserID, authInfo.Rule); err != nil {
		return errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", authInfo.Claims.Roles, authInfo.Rule, err)
	}

	return nil
}
