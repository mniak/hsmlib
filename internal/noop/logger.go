package noop

type noopLogger struct{}

func Logger() noopLogger {
	return noopLogger{}
}

func (noopLogger) Info(msg string, args ...any)  {}
func (noopLogger) Error(msg string, args ...any) {}
