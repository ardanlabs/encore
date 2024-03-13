package encore

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	encauth "encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/business/core/crud/user"
	v1 "github.com/ardanlabs/encore/business/web/v1"
	"github.com/ardanlabs/encore/business/web/v1/auth"
	"github.com/ardanlabs/encore/foundation/validate"
	"github.com/ardanlabs/encore/foundation/web"
	"github.com/google/uuid"
)

//encore:middleware target=all
func (s *Service) Context(req middleware.Request, next middleware.Next) middleware.Response {
	v := web.Values{
		TraceID: req.Data().Trace.TraceID,
		Now:     time.Now().UTC(),
	}

	req = web.SetValues(req, &v)

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
		"statuscode", resp.HTTPStatus, "since", time.Since(web.GetTime(ctx)).String())

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

	var er v1.ErrorResponse
	var status int

	switch {
	case v1.IsTrustedError(resp.Err):
		trsErr := v1.GetTrustedError(resp.Err)

		if validate.IsFieldErrors(trsErr.Err) {
			fieldErrors := validate.GetFieldErrors(trsErr.Err)
			er = v1.ErrorResponse{
				Error:  "data validation error",
				Fields: fieldErrors.Fields(),
			}
			status = trsErr.Status
			break
		}

		er = v1.ErrorResponse{
			Error: trsErr.Error(),
		}
		status = trsErr.Status

	case auth.IsAuthError(resp.Err):
		er = v1.ErrorResponse{
			Error: http.StatusText(http.StatusUnauthorized),
		}
		status = http.StatusUnauthorized

	default:
		er = v1.ErrorResponse{
			Error: http.StatusText(http.StatusInternalServerError),
		}
		status = http.StatusInternalServerError
	}

	return middleware.Response{
		HTTPStatus: status,
		Err: &errs.Error{
			Code:    errs.Internal,
			Message: "process details document",
			Details: er,
		},
	}
}

//encore:middleware target=tag:authuser
func (s *Service) AuthUser(req middleware.Request, next middleware.Next) middleware.Response {
	ctx := req.Context()
	var userID uuid.UUID

	if len(req.Data().PathParams) == 1 {
		id := req.Data().PathParams[0]
		var err error

		userID, err = uuid.Parse(id.Value)
		if err != nil {
			return middleware.Response{
				HTTPStatus: http.StatusBadRequest,
				Err:        v1.NewTrustedError(ErrInvalidID, http.StatusBadRequest),
			}
		}

		usr, err := s.usrCore.QueryByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrNotFound):
				return middleware.Response{
					HTTPStatus: http.StatusNoContent,
					Err:        v1.NewTrustedError(err, http.StatusNoContent),
				}
			default:
				return middleware.Response{
					HTTPStatus: http.StatusInternalServerError,
					Err:        fmt.Errorf("querybyid: userID[%s]: %w", userID, err),
				}
			}
		}

		ctx = setUser(ctx, usr)
	}

	claims, _ := encauth.Data().(*auth.Claims)

	if err := s.auth.Authorize(ctx, *claims, userID, auth.RuleAdminOrSubject); err != nil {
		return middleware.Response{
			HTTPStatus: http.StatusBadRequest,
			Err:        auth.NewAuthError("authorize: you are not authorized for that action, claims[%v] rule[%v]: %s", claims.Roles, auth.RuleAdminOrSubject, err),
		}
	}

	return next(req)
}
