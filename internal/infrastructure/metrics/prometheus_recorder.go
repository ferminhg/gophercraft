package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/fermin/gophercraft/internal/domain/port"
)

var _ port.MetricsRecorder = (*PrometheusRecorder)(nil)
var _ prometheus.Gatherer = (*PrometheusRecorder)(nil)

// PrometheusRecorder implements port.MetricsRecorder backed by Prometheus counters and histograms,
// registering metrics on its own Prometheus registry (not the global one).
type PrometheusRecorder struct {
	*prometheus.Registry
	prometheusMetrics
}

// NewPrometheusRecorder constructs a PrometheusRecorder with freshly registered metrics.
func NewPrometheusRecorder() *PrometheusRecorder {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &PrometheusRecorder{
		Registry:          reg,
		prometheusMetrics: newPrometheusMetrics(reg),
	}
}

// RecordHTTPRequest implements port.MetricsRecorder.
func (r *PrometheusRecorder) RecordHTTPRequest(method, route string, statusCode int, durationSecs float64) {
	if route == "" {
		route = "unknown"
	}

	r.requestsTotal.WithLabelValues(method, route, strconv.Itoa(statusCode)).Inc()
	r.requestDuration.WithLabelValues(method, route).Observe(durationSecs)
}
