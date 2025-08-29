package logger

type LoopLogger struct {
	Records []string
}

func (l *LoopLogger) Debug(msg string, args ...any) {
	l.Records = append(l.Records, "DEBUG: "+msg)
}
func (l *LoopLogger) Info(msg string, args ...any) {
	l.Records = append(l.Records, "INFO: "+msg)
}
func (l *LoopLogger) Warn(msg string, args ...any) {
	l.Records = append(l.Records, "WARN: "+msg)
}
func (l *LoopLogger) Error(msg string, args ...any) {
	l.Records = append(l.Records, "ERROR: "+msg)
}
func (l *LoopLogger) Fatal(msg string, args ...any) {
	l.Records = append(l.Records, "FATAL: "+msg)
}
func (l *LoopLogger) With(args ...any) Logger {
	return l
}
