package mid

import (
	"encore.dev/middleware"
	"github.com/ardanlabs/encore/app/sdk/metrics"
)

// Metrics sets the basic counters and guages.
func Metrics(v *metrics.Values, req middleware.Request, next middleware.Next) middleware.Response {
	n := v.IncRequests()

	if n%1000 == 0 {
		v.SetGoroutines()
	}

	resp := next(req)

	if resp.Err != nil {
		v.IncFailures()
	}

	return resp
}
