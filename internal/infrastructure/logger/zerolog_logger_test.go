package logger_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fermin/gophercraft/internal/domain/port"
	"github.com/fermin/gophercraft/internal/infrastructure/logger"
)

func TestNewZerologLoggerWithWriter_InfoIncludesKeyValues(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.NewZerologLoggerWithWriter(&buf, "debug", false, nil)

	l.Info("hello", "user_id", "42", "ok", true)

	var line map[string]any
	dec := json.NewDecoder(bytes.NewReader(bytes.TrimSpace(buf.Bytes())))
	require.NoError(t, dec.Decode(&line))

	assert.Equal(t, "hello", line["message"])
	assert.Equal(t, "42", line["user_id"])
	assert.Equal(t, true, line["ok"])
	assert.Contains(t, line, "level")
	assert.Equal(t, "info", line["level"])
}

func TestNewZerologLoggerWithWriter_WarnErrorDebugLevels(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		write     func(port.Logger)
		wantLevel string
	}{
		{
			name:      "warn",
			write:     func(l port.Logger) { l.Warn("m") },
			wantLevel: "warn",
		},
		{
			name:      "error",
			write:     func(l port.Logger) { l.Error("m") },
			wantLevel: "error",
		},
		{
			name:      "debug",
			write:     func(l port.Logger) { l.Debug("m") },
			wantLevel: "debug",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			l := logger.NewZerologLoggerWithWriter(&buf, "debug", false, nil)
			tc.write(l)

			var line map[string]any
			dec := json.NewDecoder(bytes.NewReader(bytes.TrimSpace(buf.Bytes())))
			require.NoError(t, dec.Decode(&line))
			assert.Equal(t, tc.wantLevel, line["level"])
		})
	}
}

func TestNewZerologLoggerWithWriter_WithPropagatesFields(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	base := logger.NewZerologLoggerWithWriter(&buf, "info", false, nil)
	l := base.With("trace_id", "abc")

	l.Info("done")

	var line map[string]any
	dec := json.NewDecoder(bytes.NewReader(bytes.TrimSpace(buf.Bytes())))
	require.NoError(t, dec.Decode(&line))

	assert.Equal(t, "done", line["message"])
	assert.Equal(t, "abc", line["trace_id"])
}

func TestNewZerologLoggerWithWriter_OddArgCountMarksInvalidPair(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.NewZerologLoggerWithWriter(&buf, "info", false, nil)

	l.Info("odd", "only_key")

	var line map[string]any
	dec := json.NewDecoder(bytes.NewReader(bytes.TrimSpace(buf.Bytes())))
	require.NoError(t, dec.Decode(&line))

	assert.Contains(t, line, "INVALID_LOG_PAIR")
}

func TestNewDiscardingZerologLogger_DoesNotPanic(t *testing.T) {
	t.Parallel()

	l := logger.NewDiscardingZerologLogger()
	l.Info("silenced", "k", "v")
	l.Warn("w")
	l.Error("e")
	l.Debug("d")
}

func TestNewZerologLoggerWithWriter_GlobalFieldsOnEveryLine(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.NewZerologLoggerWithWriter(&buf, "info", false, map[string]string{
		"service.name":           "demo",
		"deployment.environment": "test",
	})

	l.Info("hello")

	var line map[string]any
	dec := json.NewDecoder(bytes.NewReader(bytes.TrimSpace(buf.Bytes())))
	require.NoError(t, dec.Decode(&line))

	assert.Equal(t, "hello", line["message"])
	assert.Equal(t, "demo", line["service.name"])
	assert.Equal(t, "test", line["deployment.environment"])
}
