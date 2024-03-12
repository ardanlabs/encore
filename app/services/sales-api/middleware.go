package encore

import (
	"time"

	"encore.dev/middleware"
	"github.com/ardanlabs/encore/foundation/web"
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
	ctx := req.Context()
	er := req.Data()
	v := web.GetValues(ctx)

	// TODO: Get the query string in here.

	s.log.Info(ctx, "request started", "endpoint", er.Endpoint, "path", er.Path)

	resp := next(req)

	s.log.Info(ctx, "request completed", "endpoint", er.Endpoint, "path", er.Path,
		"statuscode", resp.HTTPStatus, "since", time.Since(v.Now).String())

	return resp
}
