// Package metrics provides infrastructure adapters for metrics recording.
package metrics

import "github.com/fermin/gophercraft/internal/domain/port"

var _ port.MetricsRecorder = (*NoopRecorder)(nil)

// NoopRecorder implements port.MetricsRecorder by discarding all samples (bootstrap / tests).
type NoopRecorder struct{}

// RecordHTTPRequest implements port.MetricsRecorder.
func (NoopRecorder) RecordHTTPRequest(string, string, int, float64) {}
