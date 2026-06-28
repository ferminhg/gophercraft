package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type prometheusMetrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func newPrometheusMetrics(reg prometheus.Registerer) prometheusMetrics {
	m := prometheusMetrics{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "route", "status_code"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "route"},
		),
	}

	reg.MustRegister(m.requestsTotal, m.requestDuration)

	return m
}
