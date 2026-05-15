package port

// Logger is a structured logging abstraction used by application and adapters.
// Methods accept alternating key-value pairs in args (like log/slog).
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Fatal(msg string, args ...any)
	With(args ...any) Logger
}
