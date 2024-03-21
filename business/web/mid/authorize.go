package mid

import (
	"context"
	"errors"
	"net/http"

	encauth "encore.dev/beta/auth"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/core/crud/user"
	"github.com/ardanlabs/encore/business/web/auth"
	"github.com/ardanlabs/encore/business/web/errs"
	"github.com/google/uuid"
)

type ctxUserKey int

const (
	userIDKey ctxUserKey = iota + 1
	userKey
)

func setUserID(req middleware.Request, userID uuid.UUID) middleware.Request {
	ctx := context.WithValue(req.Context(), userIDKey, userID)
	return req.WithContext(ctx)
}

// GetUserID extracts the user id from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found")
	}

	return v, nil
}

func setUser(req middleware.Request, usr user.User) middleware.Request {
	ctx := context.WithValue(req.Context(), userKey, usr)
	return req.WithContext(ctx)
}

// GetUser extracts the user from the context.
func GetUser(ctx context.Context) (user.User, error) {
	v, ok := ctx.Value(userKey).(user.User)
	if !ok {
		return user.User{}, errors.New("user not found")
	}

	return v, nil
}

// =============================================================================

// AuthorizeAdminOnly checks the user making the request is an admin.
func AuthorizeAdminOnly(a *auth.Auth, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	claims := encauth.Data().(*auth.Claims)

	if err := a.Authorize(ctx, *claims, uuid.UUID{}, auth.RuleAdminOnly); err != nil {
		return errs.NewResponsef(http.StatusBadRequest, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOnly, err)
	}

	return next(req)
}

// AuthorizeUser checks the user making the call has specified a user id on
// the route that matches the claims.
func AuthorizeUser(a *auth.Auth, usrCore *user.Core, req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		userID, err := uuid.Parse(id.Value)
		if err != nil {
			return errs.NewResponse(http.StatusBadRequest, ErrInvalidID)
		}

		usr, err := usrCore.QueryByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrNotFound):
				return errs.NewResponse(http.StatusBadRequest, err)

			default:
				return errs.NewResponsef(http.StatusInternalServerError, "querybyid: userID[%s]: %s", userID, err)
			}
		}

		req = setUser(req, usr)
	}

	claims := encauth.Data().(*auth.Claims)

	if err := a.Authorize(ctx, *claims, userID, auth.RuleAdminOrSubject); err != nil {
		return errs.NewResponsef(http.StatusBadRequest, "authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err)
	}

	return next(req)
}
