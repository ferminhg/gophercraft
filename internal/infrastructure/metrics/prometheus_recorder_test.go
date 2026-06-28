package metrics

import (
	"net/http"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dtomodel "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrometheusRecorder_RecordHTTPRequest(t *testing.T) {
	t.Parallel()

	recorder := NewPrometheusRecorder()
	recorder.RecordHTTPRequest(http.MethodGet, "/status", http.StatusOK, 0.05)

	counter, err := recorder.requestsTotal.GetMetricWithLabelValues(http.MethodGet, "/status", "200")
	require.NoError(t, err)
	assert.Equal(t, 1.0, testutil.ToFloat64(counter))

	sum := findHistogramSampleSum(t, recorder, "http_request_duration_seconds", map[string]string{
		"method": http.MethodGet,
		"route":  "/status",
	})
	assert.InDelta(t, 0.05, sum, 1e-9)
}

func TestPrometheusRecorder_IncludesGoAndProcessMetrics(t *testing.T) {
	t.Parallel()

	recorder := NewPrometheusRecorder()
	mfs, err := recorder.Gather()
	require.NoError(t, err)

	names := make(map[string]bool, len(mfs))
	for _, mf := range mfs {
		names[mf.GetName()] = true
	}

	assert.True(t, names["go_goroutines"])
	assert.True(t, names["process_cpu_seconds_total"])
}

func TestPrometheusRecorder_RecordHTTPRequest_EmptyRouteUsesUnknown(t *testing.T) {
	t.Parallel()

	recorder := NewPrometheusRecorder()
	recorder.RecordHTTPRequest(http.MethodPost, "", http.StatusNotFound, 0)

	counter, err := recorder.requestsTotal.GetMetricWithLabelValues(http.MethodPost, "unknown", "404")
	require.NoError(t, err)
	assert.Equal(t, 1.0, testutil.ToFloat64(counter))
}

func findHistogramSampleSum(t *testing.T, g prometheus.Gatherer, metricName string, wantLabels map[string]string) float64 {
	t.Helper()

	mfs, err := g.Gather()
	require.NoError(t, err)

	for _, mf := range mfs {
		if mf.GetName() != metricName {
			continue
		}
		for _, m := range mf.GetMetric() {
			if metricLabelsMatch(m, wantLabels) {
				require.NotNil(t, m.Histogram)
				return m.Histogram.GetSampleSum()
			}
		}
	}

	t.Fatalf("metric %q with labels %v not found", metricName, wantLabels)
	return 0
}

func metricLabelsMatch(m *dtomodel.Metric, want map[string]string) bool {
	got := make(map[string]string, len(m.Label))
	for _, lp := range m.Label {
		got[lp.GetName()] = lp.GetValue()
	}
	for k, v := range want {
		if got[k] != v {
			return false
		}
	}
	return true
}
