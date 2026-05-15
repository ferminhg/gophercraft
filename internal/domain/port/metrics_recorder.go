package port

// MetricsRecorder records outbound observability signals (e.g. Prometheus metrics).
type MetricsRecorder interface {
	RecordHTTPRequest(method, route string, statusCode int, durationSecs float64)
}
