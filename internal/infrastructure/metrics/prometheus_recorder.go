package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/fermin/gophercraft/internal/domain/port"
)

var _ port.MetricsRecorder = (*PrometheusRecorder)(nil)
var _ prometheus.Gatherer = (*PrometheusRecorder)(nil)

// PrometheusRecorder implements port.MetricsRecorder backed by Prometheus counters and histograms,
// registering metrics on its own Prometheus registry (not the global one).
type PrometheusRecorder struct {
	*prometheus.Registry

	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

// NewPrometheusRecorder constructs a PrometheusRecorder with freshly registered metrics.
func NewPrometheusRecorder() *PrometheusRecorder {
	reg := prometheus.NewRegistry()

	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status_code"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	reg.MustRegister(requestsTotal, requestDuration)

	return &PrometheusRecorder{
		Registry:        reg,
		requestsTotal:   requestsTotal,
		requestDuration: requestDuration,
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
