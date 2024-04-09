// Package metrics provides support to update metric values.
package metrics

import (
	"expvar"
	"runtime"

	"encore.dev"
	"encore.dev/metrics"
)

var devGoroutines = expvar.NewInt("goroutines")
var devRequests = expvar.NewInt("requests")
var devFailures = expvar.NewInt("errors")
var devPanics = expvar.NewInt("panics")

// Config lists the set of metrics that is tracked.
type Config struct {
	Goroutines *metrics.Gauge[uint64]
	Requests   *metrics.Counter[uint64]
	Failures   *metrics.Counter[uint64]
	Panics     *metrics.Counter[uint64]
}

// Values provides an api to work with metrics.
type Values struct {
	devEnv        bool
	goroutines    *metrics.Gauge[uint64]
	requests      *metrics.Counter[uint64]
	failures      *metrics.Counter[uint64]
	panics        *metrics.Counter[uint64]
	devGoroutines *expvar.Int
	devRequests   *expvar.Int
	devFailures   *expvar.Int
	devPanics     *expvar.Int
}

// New constructs a Values for working with metrics.
func New(cfg Config) *Values {
	return &Values{
		devEnv:        encore.Meta().Environment.Type == encore.EnvDevelopment,
		goroutines:    cfg.Goroutines,
		requests:      cfg.Requests,
		failures:      cfg.Failures,
		panics:        cfg.Panics,
		devGoroutines: devGoroutines,
		devRequests:   devRequests,
		devFailures:   devFailures,
		devPanics:     devPanics,
	}
}

// SetGoroutines updates the number of goroutines.
func (v *Values) SetGoroutines() {
	n := runtime.NumGoroutine()
	v.goroutines.Set(uint64(n))

	if v.devEnv {
		v.devGoroutines.Set(int64(n))
	}
}

// IncRequests increments the requests by 1.
func (v *Values) IncRequests() int64 {
	v.requests.Add(1)
	v.devRequests.Add(1)

	return v.devRequests.Value()
}

// IncFailures increments the failures by 1.
func (v *Values) IncFailures() {
	v.failures.Add(1)

	if v.devEnv {
		v.devFailures.Add(1)
	}
}

// IncPanics increments the panics by 1.
func (v *Values) IncPanics() {
	v.panics.Add(1)

	if v.devEnv {
		v.devPanics.Add(1)
	}
}
