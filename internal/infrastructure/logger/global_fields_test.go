package logger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fermin/gophercraft/internal/infrastructure/logger"
)

func TestGlobalFieldsFromEnv_OTELServiceName(t *testing.T) {
	t.Setenv("OTEL_SERVICE_NAME", "my-api")
	t.Setenv("SERVICE_NAME", "ignored-when-otel-set")
	t.Setenv("DEPLOYMENT_ENVIRONMENT", "")
	t.Setenv("ENV", "")
	t.Setenv("APP_ENV", "")

	fields := logger.GlobalFieldsFromEnv()
	require.Contains(t, fields, "service.name")
	assert.Equal(t, "my-api", fields["service.name"])
	assert.NotContains(t, fields, "deployment.environment")
}

func TestGlobalFieldsFromEnv_FallbacksAndDeployment(t *testing.T) {
	t.Setenv("OTEL_SERVICE_NAME", "")
	t.Setenv("SERVICE_NAME", "from-service")
	t.Setenv("DEPLOYMENT_ENVIRONMENT", "staging")

	fields := logger.GlobalFieldsFromEnv()
	assert.Equal(t, "from-service", fields["service.name"])
	assert.Equal(t, "staging", fields["deployment.environment"])
}

func TestGlobalFieldsFromEnv_DefaultServiceName(t *testing.T) {
	t.Setenv("OTEL_SERVICE_NAME", "")
	t.Setenv("SERVICE_NAME", "")
	t.Setenv("DEPLOYMENT_ENVIRONMENT", "")
	t.Setenv("ENV", "")
	t.Setenv("APP_ENV", "")

	fields := logger.GlobalFieldsFromEnv()
	assert.Equal(t, "gophercraft", fields["service.name"])
}

func TestGlobalFieldsFromEnv_AppEnvFallback(t *testing.T) {
	t.Setenv("OTEL_SERVICE_NAME", "x")
	t.Setenv("DEPLOYMENT_ENVIRONMENT", "")
	t.Setenv("ENV", "")
	t.Setenv("APP_ENV", "dev")

	fields := logger.GlobalFieldsFromEnv()
	assert.Equal(t, "dev", fields["deployment.environment"])
}
