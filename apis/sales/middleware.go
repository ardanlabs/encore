package sales

import (
	"fmt"

	eerrs "encore.dev/beta/errs"
	"encore.dev/middleware"
	authsrv "github.com/ardanlabs/encore/apis/auth"
	"github.com/ardanlabs/encore/app/api/errs"
	"github.com/ardanlabs/encore/app/api/mid"
)

// =============================================================================
// Global middleware functions
// The order matters so be careful when injecting new middleware.

//lint:ignore U1000 "called by encore"
//encore:middleware target=all
func (s *Service) panics(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.Panics(s.mtrcs, req, next)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:metrics
func (s *Service) metrics(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.Metrics(s.mtrcs, req, next)
}

// =============================================================================
// Authorization related middleware
// These middleware functions must come after the global middleware functions
// above. These are targeted so the order doesn't matter.

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize
func (s *Service) authorize(req middleware.Request, next middleware.Next) middleware.Response {
	p, req, err := mid.Authorize(req, next)
	if err != nil {
		return errs.NewResponse(eerrs.Unauthenticated, err)
	}

	if err := authsrv.Authorize(req.Context(), p); err != nil {
		err = fmt.Errorf("%s", err.Error()[17:]) // Remove "unauthenticated:" from the error string.
		return errs.NewResponse(eerrs.Unauthenticated, err)
	}

	return next(req)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_user
func (s *Service) authorizeUser(req middleware.Request, next middleware.Next) middleware.Response {
	p, req, err := mid.AuthorizeUser(s.userBus, req, next)
	if err != nil {
		return errs.NewResponse(eerrs.Unauthenticated, err)
	}

	if err := authsrv.Authorize(req.Context(), p); err != nil {
		err = fmt.Errorf("%s", err.Error()[17:]) // Remove "unauthenticated:" from the error string.
		return errs.NewResponse(eerrs.Unauthenticated, err)
	}

	return next(req)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_product
func (s *Service) authorizeProduct(req middleware.Request, next middleware.Next) middleware.Response {
	p, req, err := mid.AuthorizeProduct(s.productBus, req, next)
	if err != nil {
		return errs.NewResponse(eerrs.Unauthenticated, err)
	}

	if err := authsrv.Authorize(req.Context(), p); err != nil {
		err = fmt.Errorf("%s", err.Error()[17:]) // Remove "unauthenticated:" from the error string.
		return errs.NewResponse(eerrs.Unauthenticated, err)
	}

	return next(req)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_home
func (s *Service) authorizeHome(req middleware.Request, next middleware.Next) middleware.Response {
	p, req, err := mid.AuthorizeHome(s.homeBus, req, next)
	if err != nil {
		return errs.NewResponse(eerrs.Unauthenticated, err)
	}

	if err := authsrv.Authorize(req.Context(), p); err != nil {
		err = fmt.Errorf("%s", err.Error()[17:]) // Remove "unauthenticated:" from the error string.
		return errs.NewResponse(eerrs.Unauthenticated, err)
	}

	return next(req)
}
