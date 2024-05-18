package sales

import (
	"context"
	"fmt"
	"time"

	"encore.dev/middleware"
	authsrv "github.com/ardanlabs/encore/api/services/auth"
	"github.com/ardanlabs/encore/app/sdk/errs"
	"github.com/ardanlabs/encore/app/sdk/mid"
	"github.com/ardanlabs/encore/business/sdk/sqldb"
)

// NOTE: The order matters so be careful when injecting new middleware. Global
//       middleware will always come first. We want the Auth middleware to
//       happen before any non-global middlware.

// =============================================================================
// Global middleware functions

//lint:ignore U1000 "called by encore"
//encore:middleware target=all
func (s *Service) panics(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.Panics(s.mtrcs, req, next)
}

// =============================================================================
// Authorization related middleware

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize
func (s *Service) authorize(req middleware.Request, next middleware.Next) middleware.Response {
	p, req, err := mid.Authorize(req)
	if err != nil {
		return errs.NewResponse(errs.Unauthenticated, err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	if err := authsrv.Authorize(ctx, p); err != nil {
		err = fmt.Errorf("%s", err.Error()[17:]) // Remove "unauthenticated:" from the error string.
		return errs.NewResponse(errs.Unauthenticated, err)
	}

	return next(req)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_user
func (s *Service) authorizeUser(req middleware.Request, next middleware.Next) middleware.Response {
	p, req, err := mid.AuthorizeUser(s.userBus, req)
	if err != nil {
		return errs.NewResponse(errs.Unauthenticated, err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	if err := authsrv.Authorize(ctx, p); err != nil {
		err = fmt.Errorf("%s", err.Error()[17:]) // Remove "unauthenticated:" from the error string.
		return errs.NewResponse(errs.Unauthenticated, err)
	}

	return next(req)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_product
func (s *Service) authorizeProduct(req middleware.Request, next middleware.Next) middleware.Response {
	p, req, err := mid.AuthorizeProduct(s.productBus, req)
	if err != nil {
		return errs.NewResponse(errs.Unauthenticated, err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	if err := authsrv.Authorize(ctx, p); err != nil {
		err = fmt.Errorf("%s", err.Error()[17:]) // Remove "unauthenticated:" from the error string.
		return errs.NewResponse(errs.Unauthenticated, err)
	}

	return next(req)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:authorize_home
func (s *Service) authorizeHome(req middleware.Request, next middleware.Next) middleware.Response {
	p, req, err := mid.AuthorizeHome(s.homeBus, req)
	if err != nil {
		return errs.NewResponse(errs.Unauthenticated, err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	if err := authsrv.Authorize(ctx, p); err != nil {
		err = fmt.Errorf("%s", err.Error()[17:]) // Remove "unauthenticated:" from the error string.
		return errs.NewResponse(errs.Unauthenticated, err)
	}

	return next(req)
}

// =============================================================================
// Specific middleware functions

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:transaction
func (s *Service) beginCommitRollback(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.BeginCommitRollback(s.log, sqldb.NewBeginner(s.db), req, next)
}

//lint:ignore U1000 "called by encore"
//encore:middleware target=tag:metrics
func (s *Service) metrics(req middleware.Request, next middleware.Next) middleware.Response {
	return mid.Metrics(s.mtrcs, req, next)
}
