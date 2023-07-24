package noop

import "github.com/mniak/hsmlib"

type noopLogger struct{}

func Logger() hsmlib.Logger {
	return noopLogger{}
}

func (noopLogger) Info(msg string, args ...any)  {}
func (noopLogger) Error(msg string, args ...any) {}
