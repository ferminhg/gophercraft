package logger

import "os"

// GlobalFieldsFromEnv returns OpenTelemetry-style resource attributes for every log line
// (see https://opentelemetry.io/docs/specs/semconv/resource/#service and deployment).
func GlobalFieldsFromEnv() map[string]string {
	svc := os.Getenv("OTEL_SERVICE_NAME")
	if svc == "" {
		svc = os.Getenv("SERVICE_NAME")
	}
	if svc == "" {
		svc = "gophercraft"
	}

	m := map[string]string{"service.name": svc}

	env := os.Getenv("DEPLOYMENT_ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	if env != "" {
		m["deployment.environment"] = env
	}

	return m
}
