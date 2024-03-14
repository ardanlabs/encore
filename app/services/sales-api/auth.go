package encore

import (
	"context"
	"errors"
	"net/http"

	encauth "encore.dev/beta/auth"
	"encore.dev/beta/errs"
	v1 "github.com/ardanlabs/encore/business/web/v1"
	"github.com/ardanlabs/encore/business/web/v1/auth"
	"github.com/google/uuid"
)

// Set of error variables for handling auth errors.
var (
	ErrInvalidID = errors.New("ID is not in its proper form")
)

//encore:authhandler
func (s *Service) AuthHandler(ctx context.Context, token string) (encauth.UID, *auth.Claims, error) {
	claims, err := s.auth.Authenticate(ctx, "Bearer "+token)
	if err != nil {
		s.log.Error(ctx, "authenticate: failed", "ERROR", err)
		return "", nil, v1.NewError(errs.Unauthenticated, http.StatusText(http.StatusUnauthorized))
	}

	if claims.Subject == "" {
		s.log.Error(ctx, "authorize: you are not authorized for that action, no claims")
		return "", nil, v1.NewError(errs.Unauthenticated, http.StatusText(http.StatusUnauthorized))
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		s.log.Error(ctx, "parsing subject: %s", ErrInvalidID)
		return "", nil, v1.NewError(errs.InvalidArgument, ErrInvalidID.Error())
	}

	return encauth.UID(subjectID.String()), &claims, nil
}
