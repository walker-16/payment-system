package logger

import (
	"io"
	"log/slog"
)

type Format string

const (
	FormatText Format = "TEXT"
	FormatJSON Format = "JSON"
)

type LoggerConfig struct {
	Format    Format
	Level     slog.Level
	Output    io.Writer
	AddSource bool
}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	With(args ...any) Logger
}

func String(k, s string) any {
	return slog.String(k, s)
}

func Error(err error) any {
	return slog.Any("error", err)
}

func Int(k string, v int) any {
	return slog.Int(k, v)
}
