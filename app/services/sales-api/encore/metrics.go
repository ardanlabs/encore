package encore

import (
	emetrics "encore.dev/metrics"
	"github.com/ardanlabs/encore/business/web/metrics"
)

// These are the counters and guages we are tracking in the app.
var goroutines = emetrics.NewGauge[uint64]("goroutines", emetrics.GaugeConfig{})
var requests = emetrics.NewCounter[uint64]("requests", emetrics.CounterConfig{})
var failures = emetrics.NewCounter[uint64]("errors", emetrics.CounterConfig{})
var panics = emetrics.NewCounter[uint64]("panics", emetrics.CounterConfig{})

func newMetrics() *metrics.Values {
	return metrics.New(metrics.Config{
		Goroutines: goroutines,
		Requests:   requests,
		Failures:   failures,
		Panics:     panics,
	})
}
