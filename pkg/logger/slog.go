package logger

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	base *slog.Logger
}

func NewSlogLogger(cfg LoggerConfig) Logger {
	output := cfg.Output
	if output == nil {
		output = os.Stdout
	}

	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler
	switch cfg.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(output, opts)
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	return &SlogLogger{slog.New(handler)}
}

func (s *SlogLogger) Debug(msg string, args ...any) {
	s.base.Debug(msg, args...)
}

func (s *SlogLogger) Info(msg string, args ...any) {
	s.base.Info(msg, args...)
}

func (s *SlogLogger) Warn(msg string, args ...any) {
	s.base.Warn(msg, args...)
}

func (s *SlogLogger) Error(msg string, args ...any) {
	s.base.Error(msg, args...)
}

func (s *SlogLogger) Fatal(msg string, args ...any) {
	s.base.Error(msg, args...)
	os.Exit(1)
}

func (s *SlogLogger) With(args ...any) Logger {
	return &SlogLogger{base: s.base.With(args...)}
}
