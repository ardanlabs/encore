package encore

import (
	emetrics "encore.dev/metrics"
	"github.com/ardanlabs/encore/business/web/metrics"
)

//lint:ignore U1000 "used by encore"
var (
	goroutines = emetrics.NewGauge[uint64]("goroutines", emetrics.GaugeConfig{})
	requests   = emetrics.NewCounter[uint64]("requests", emetrics.CounterConfig{})
	failures   = emetrics.NewCounter[uint64]("errors", emetrics.CounterConfig{})
	panics     = emetrics.NewCounter[uint64]("panics", emetrics.CounterConfig{})
)

func newMetrics() *metrics.Values {
	return metrics.New(metrics.Config{
		Goroutines: goroutines,
		Requests:   requests,
		Failures:   failures,
		Panics:     panics,
	})
}
