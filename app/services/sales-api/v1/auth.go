package encore

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	encauth "encore.dev/beta/auth"
	"encore.dev/beta/errs"
	v1 "github.com/ardanlabs/encore/business/web/v1"
	"github.com/ardanlabs/encore/business/web/v1/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var errInvalidID = errors.New("ID is not in its proper form")

type authParams struct {
	Authorization string `header:"Authorization"`
}

//encore:authhandler
func (s *Service) authHandler(ctx context.Context, ap *authParams) (encauth.UID, *auth.Claims, error) {
	parts := strings.Split(ap.Authorization, " ")
	if len(parts) != 2 {
		return "", nil, v1.NewError(errs.Unauthenticated, "invalid authorization value")
	}

	switch parts[0] {
	case "Bearer":
		return s.processJWT(ctx, ap.Authorization)

	case "Basic":
		return s.processBasic(ctx, ap.Authorization)
	}

	return "", nil, v1.NewError(errs.Unauthenticated, http.StatusText(http.StatusUnauthorized))
}

func (s *Service) processJWT(ctx context.Context, token string) (encauth.UID, *auth.Claims, error) {
	claims, err := s.auth.Authenticate(ctx, token)
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
		s.log.Error(ctx, "parsing subject: %s", errInvalidID)
		return "", nil, v1.NewError(errs.InvalidArgument, errInvalidID.Error())
	}

	return encauth.UID(subjectID.String()), &claims, nil

}

func (s *Service) processBasic(ctx context.Context, basic string) (encauth.UID, *auth.Claims, error) {
	email, pass, ok := parseBasicAuth(basic)
	if !ok {
		return "", nil, v1.NewError(errs.Unauthenticated, "invalid Basic auth")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", nil, v1.NewError(errs.Unauthenticated, "invalid email format")
	}

	usr, err := s.usrCore.Authenticate(ctx, *addr, pass)
	if err != nil {
		return "", nil, v1.NewError(errs.Unauthenticated, err.Error())
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		s.log.Error(ctx, "parsing subject: %s", errInvalidID)
		return "", nil, v1.NewError(errs.InvalidArgument, errInvalidID.Error())
	}

	return encauth.UID(subjectID.String()), &claims, nil
}

func parseBasicAuth(auth string) (string, string, bool) {
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", false
	}

	c, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}

	username, password, ok := strings.Cut(string(c), ":")
	if !ok {
		return "", "", false
	}

	return username, password, true
}
