package mid

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	encauth "encore.dev/beta/auth"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/web/auth"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// AuthParams is used to unmarshal the authorization string from the request.
type AuthParams struct {
	Authorization string `header:"Authorization"`
}

// =============================================================================

// AuthHandler is used to provide initial auth for JWT's and basic user:password.
func AuthHandler(ctx context.Context, a *auth.Auth, usrCore *user.Core, ap *AuthParams) (encauth.UID, *auth.Claims, error) {
	parts := strings.Split(ap.Authorization, " ")
	if len(parts) != 2 {
		return "", nil, errs.Newf(http.StatusUnauthorized, "invalid authorization value")
	}

	switch parts[0] {
	case "Bearer":
		return processJWT(ctx, a, ap.Authorization)

	case "Basic":
		return processBasic(ctx, usrCore, ap.Authorization)
	}

	return "", nil, errs.Newf(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}

// =============================================================================

func processJWT(ctx context.Context, a *auth.Auth, token string) (encauth.UID, *auth.Claims, error) {
	claims, err := a.Authenticate(ctx, token)
	if err != nil {
		return "", nil, errs.New(http.StatusUnauthorized, err)
	}

	if claims.Subject == "" {
		return "", nil, errs.Newf(http.StatusUnauthorized, "authorize: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", nil, errs.New(http.StatusUnauthorized, fmt.Errorf("parsing subject: %w", err))
	}

	return encauth.UID(subjectID.String()), &claims, nil
}

func processBasic(ctx context.Context, usrCore *user.Core, basic string) (encauth.UID, *auth.Claims, error) {
	email, pass, ok := parseBasicAuth(basic)
	if !ok {
		return "", nil, errs.Newf(http.StatusUnauthorized, "invalid Basic auth")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", nil, errs.New(http.StatusUnauthorized, err)
	}

	usr, err := usrCore.Authenticate(ctx, *addr, pass)
	if err != nil {
		return "", nil, errs.New(http.StatusUnauthorized, err)
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
		return "", nil, errs.Newf(http.StatusUnauthorized, "parsing subject: %s", err)
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
