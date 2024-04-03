package salesapiweb

import (
	emetrics "encore.dev/metrics"
	"github.com/ardanlabs/encore/app/api/metrics"
)

// Encore currently requires these metrics to be declared in the same package
// as the service type.
//
//lint:ignore U1000 "used by encore"
var (
	goroutines = emetrics.NewGauge[uint64]("goroutines", emetrics.GaugeConfig{})
	requests   = emetrics.NewCounter[uint64]("requests", emetrics.CounterConfig{})
	failures   = emetrics.NewCounter[uint64]("errors", emetrics.CounterConfig{})
	panics     = emetrics.NewCounter[uint64]("panics", emetrics.CounterConfig{})
)

// newMetrics will construct a business layer metrics value that will allow
// the metrics above to be passed to the business layer metrics middleware
// function. Remember, business layer packages can't import app layer packages.
func newMetrics() *metrics.Values {
	return metrics.New(metrics.Config{
		Goroutines: goroutines,
		Requests:   requests,
		Failures:   failures,
		Panics:     panics,
	})
}
