package logger

type NoopLogger struct{}

func NewNoopLogger() Logger {
	return &NoopLogger{}
}

func (n *NoopLogger) Debug(msg string, args ...any) {
}

func (n *NoopLogger) Info(msg string, args ...any) {
}

func (n *NoopLogger) Warn(msg string, args ...any) {
}

func (n *NoopLogger) Error(msg string, args ...any) {
}

func (n *NoopLogger) Fatal(msg string, args ...any) {
}

func (n *NoopLogger) With(args ...any) Logger {
	return n
}
