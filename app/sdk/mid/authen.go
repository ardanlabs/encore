package mid

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/mail"
	"strings"
	"time"

	eauth "encore.dev/beta/auth"
	eerrs "encore.dev/beta/errs"
	"github.com/ardanlabs/encore/app/sdk/auth"
	"github.com/ardanlabs/encore/app/sdk/errs"
	"github.com/ardanlabs/encore/business/domain/userbus"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Bearer processes JWT authentication logic.
func Bearer(ctx context.Context, ath *auth.Auth, authorization string) (eauth.UID, *auth.Claims, error) {
	claims, err := ath.Authenticate(ctx, authorization)
	if err != nil {
		return "", nil, errs.New(eerrs.Unauthenticated, err)
	}

	if claims.Subject == "" {
		return "", nil, errs.Newf(eerrs.Unauthenticated, "authorize: you are not authorized for that action, no claims")
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", nil, errs.New(eerrs.Unauthenticated, fmt.Errorf("parsing subject: %w", err))
	}

	return eauth.UID(subjectID.String()), &claims, nil
}

// Basic processes basic authentication logic.
func Basic(ctx context.Context, ath *auth.Auth, userBus *userbus.Business, authorization string) (eauth.UID, *auth.Claims, error) {
	email, pass, ok := parseBasicAuth(authorization)
	if !ok {
		return "", nil, errs.Newf(eerrs.Unauthenticated, "invalid Basic auth")
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", nil, errs.New(eerrs.Unauthenticated, err)
	}

	usr, err := userBus.Authenticate(ctx, *addr, pass)
	if err != nil {
		return "", nil, errs.New(eerrs.Unauthenticated, err)
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    ath.Issuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: usr.Roles,
	}

	subjectID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return "", nil, errs.Newf(eerrs.Unauthenticated, "parsing subject: %s", err)
	}

	return eauth.UID(subjectID.String()), &claims, nil
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
