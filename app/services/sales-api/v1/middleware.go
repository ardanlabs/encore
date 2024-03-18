package encore

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	encauth "encore.dev/beta/auth"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/core/crud/user"
	v1 "github.com/ardanlabs/encore/business/web/v1"
	"github.com/ardanlabs/encore/business/web/v1/auth"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/uuid"
)

//encore:middleware target=all
func (s *Service) Context(req middleware.Request, next middleware.Next) middleware.Response {
	v := values{
		TraceID: req.Data().Trace.TraceID,
		Now:     time.Now().UTC(),
	}

	req = setValues(req, &v)

	return next(req)
}

//encore:middleware target=all
func (s *Service) Logger(req middleware.Request, next middleware.Next) middleware.Response {
	type queryParameters interface {
		Params() map[string]string
	}

	ctx := req.Context()
	er := req.Data()

	path := er.Path
	if qp, ok := er.Payload.(queryParameters); ok {
		var b strings.Builder
		fmt.Fprintf(&b, "%s?", path)
		for k, v := range qp.Params() {
			fmt.Fprintf(&b, "%s=%s&", k, v)
		}
		path = b.String()
		path = path[:len(path)-1]
	}

	s.log.Info(ctx, "request started", "endpoint", er.Endpoint, "path", path)

	resp := next(req)

	s.log.Info(ctx, "request completed", "endpoint", er.Endpoint, "path", path,
		"statuscode", resp.HTTPStatus, "since", time.Since(getTime(ctx)).String())

	return resp
}

//encore:middleware target=all
func (s *Service) Errors(req middleware.Request, next middleware.Next) middleware.Response {
	resp := next(req)
	if resp.Err == nil {
		return resp
	}

	ctx := req.Context()

	s.log.Error(ctx, "errors message", "msg", resp.Err)

	var midResp middleware.Response

	switch {
	case v1.IsTrustedError(resp.Err):
		trsErr := v1.GetTrustedError(resp.Err)

		if validate.IsFieldErrors(trsErr.Err) {
			fieldErrors := validate.GetFieldErrors(trsErr.Err)
			midResp = v1.NewErrorResponseWithFields(trsErr.Status, "data validation error", fieldErrors)
			break
		}

		midResp = v1.NewErrorResponse(trsErr.Status, trsErr)

	case auth.IsAuthError(resp.Err):
		midResp = v1.NewErrorResponseWithMessage(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))

	default:
		midResp = v1.NewErrorResponseWithMessage(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return midResp
}

//encore:middleware target=tag:auth_admin_only
func (s *Service) AuthorizeAdminOnly(req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	claims := encauth.Data().(*auth.Claims)

	if err := s.auth.Authorize(ctx, *claims, uuid.UUID{}, auth.RuleAdminOnly); err != nil {
		authErr := auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOnly, err)
		return v1.NewErrorResponse(http.StatusBadRequest, authErr)
	}

	return next(req)
}

//encore:middleware target=tag:auth_user
func (s *Service) AuthUser(req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		userID, err := uuid.Parse(id.Value)
		if err != nil {
			return v1.NewErrorResponse(http.StatusBadRequest, ErrInvalidID)
		}

		usr, err := s.usrCore.QueryByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrNotFound):
				return v1.NewErrorResponse(http.StatusNoContent, err)

			default:
				return v1.NewErrorResponse(http.StatusInternalServerError, fmt.Errorf("querybyid: userID[%s]: %w", userID, err))
			}
		}

		ctx = setUser(ctx, usr)
	}

	claims := encauth.Data().(*auth.Claims)

	if err := s.auth.Authorize(ctx, *claims, userID, auth.RuleAdminOrSubject); err != nil {
		authErr := auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err)
		return v1.NewErrorResponse(http.StatusBadRequest, authErr)
	}

	return next(req)
}
