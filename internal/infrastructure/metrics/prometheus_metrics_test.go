package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrometheusMetrics_RegistersHTTPMetrics(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	metrics := newPrometheusMetrics(reg)

	require.NotNil(t, metrics.requestsTotal)
	require.NotNil(t, metrics.requestDuration)

	metrics.requestsTotal.WithLabelValues("GET", "/test", "200").Inc()
	metrics.requestDuration.WithLabelValues("GET", "/test").Observe(0.01)

	mfs, err := reg.Gather()
	require.NoError(t, err)

	names := make(map[string]bool, len(mfs))
	for _, mf := range mfs {
		names[mf.GetName()] = true
	}

	assert.True(t, names["http_requests_total"])
	assert.True(t, names["http_request_duration_seconds"])
}
