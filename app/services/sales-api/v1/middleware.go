package encore

import (
	"errors"
	"fmt"
	"net/http"

	encauth "encore.dev/beta/auth"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/services/sales-api/v1/handlers/usergrp"
	"github.com/ardanlabs/encore/business/core/crud/user"
	v1 "github.com/ardanlabs/encore/business/web/v1"
	"github.com/ardanlabs/encore/business/web/v1/auth"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/google/uuid"
)

// =============================================================================
// Global middleware functions.
// The order matters so be careful when injecting new middleware.

//encore:middleware target=all
func (s *Service) context(req middleware.Request, next middleware.Next) middleware.Response {
	req = setTraceID(req, req.Data().Trace.TraceID)

	return next(req)
}

//encore:middleware target=all
func (s *Service) errors(req middleware.Request, next middleware.Next) middleware.Response {
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

// =============================================================================
// Auth related middleware
// These middleware functions must come after the global middleware functions
// above. These are targeted so the order doesn't matter.

//encore:middleware target=tag:auth_admin_only
func (s *Service) authorizeAdminOnly(req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	claims := encauth.Data().(*auth.Claims)

	if err := s.auth.Authorize(ctx, *claims, uuid.UUID{}, auth.RuleAdminOnly); err != nil {
		authErr := auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOnly, err)
		return v1.NewErrorResponse(http.StatusBadRequest, authErr)
	}

	return next(req)
}

//encore:middleware target=tag:auth_user
func (s *Service) authUser(req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]

		userID, err := uuid.Parse(id.Value)
		if err != nil {
			return v1.NewErrorResponse(http.StatusBadRequest, errInvalidID)
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

		req = usergrp.SetUser(req, usr)
	}

	claims := encauth.Data().(*auth.Claims)

	if err := s.auth.Authorize(ctx, *claims, userID, auth.RuleAdminOrSubject); err != nil {
		authErr := auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err)
		return v1.NewErrorResponse(http.StatusBadRequest, authErr)
	}

	return next(req)
}
