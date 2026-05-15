// Package logger provides infrastructure adapters for structured logging.
package logger

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/rs/zerolog"

	"github.com/fermin/gophercraft/internal/domain/port"
)

var _ port.Logger = (*zerologLogger)(nil)

type zerologLogger struct {
	logger zerolog.Logger
}

// NewZerologLogger builds a structured [port.Logger] writing to stdout.
// level is parsed with zerolog (e.g. "debug", "info"); invalid values default to info.
// When pretty is true, logs use a human-readable console format; otherwise JSON lines.
// global may be nil; keys should use stable names (e.g. OpenTelemetry resource attributes
// like service.name, deployment.environment) so aggregators can correlate logs.
func NewZerologLogger(level string, pretty bool, global map[string]string) port.Logger {
	return NewZerologLoggerWithWriter(os.Stdout, level, pretty, global)
}

// NewZerologLoggerWithWriter is like [NewZerologLogger] but writes to w (tests, DI).
// global may be nil.
func NewZerologLoggerWithWriter(w io.Writer, level string, pretty bool, global map[string]string) port.Logger {
	out := w
	if pretty {
		out = zerolog.ConsoleWriter{
			Out:        w,
			TimeFormat: time.RFC3339,
		}
	}

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	ctx := zerolog.New(out).
		Level(lvl).
		With().
		Timestamp()

	ctx = applyGlobalFields(ctx, global)
	zl := ctx.Logger()

	return &zerologLogger{logger: zl}
}

func applyGlobalFields(ctx zerolog.Context, global map[string]string) zerolog.Context {
	if len(global) == 0 {
		return ctx
	}

	keys := make([]string, 0, len(global))
	for k := range global {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := global[k]
		if v == "" {
			continue
		}
		ctx = ctx.Str(k, v)
	}

	return ctx
}

// NewDiscardingZerologLogger returns a logger that discards output (e.g. HTTP tests).
func NewDiscardingZerologLogger() port.Logger {
	return NewZerologLoggerWithWriter(io.Discard, "debug", false, nil)
}

func (z *zerologLogger) Info(msg string, args ...any) {
	applyPairs(z.logger.Info(), args).Msg(msg)
}

func (z *zerologLogger) Warn(msg string, args ...any) {
	applyPairs(z.logger.Warn(), args).Msg(msg)
}

func (z *zerologLogger) Error(msg string, args ...any) {
	applyPairs(z.logger.Error(), args).Msg(msg)
}

func (z *zerologLogger) Debug(msg string, args ...any) {
	applyPairs(z.logger.Debug(), args).Msg(msg)
}

func (z *zerologLogger) Fatal(msg string, args ...any) {
	applyPairs(z.logger.Fatal(), args).Msg(msg)
}

func (z *zerologLogger) With(args ...any) port.Logger {
	ctx := z.logger.With()
	n := len(args)
	for i := 0; i < n; i += 2 {
		if i+1 >= n {
			ctx = ctx.Interface("INVALID_LOG_PAIR", args[i])
			continue
		}
		ctx = ctx.Interface(keyString(args[i]), args[i+1])
	}

	return &zerologLogger{logger: ctx.Logger()}
}

func applyPairs(e *zerolog.Event, args []any) *zerolog.Event {
	n := len(args)
	for i := 0; i < n; i += 2 {
		if i+1 >= n {
			return e.Interface("INVALID_LOG_PAIR", args[i])
		}
		e = e.Interface(keyString(args[i]), args[i+1])
	}

	return e
}

func keyString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}

	return fmt.Sprint(v)
}
