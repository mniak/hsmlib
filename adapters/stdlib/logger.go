package stdlib

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type stdlogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger(prefix string) stdlogger {
	return stdlogger{
		infoLogger:  log.New(os.Stdout, prefix+"INFO ", 0),
		errorLogger: log.New(os.Stderr, prefix+"ERROR ", 0),
	}
}

func (l stdlogger) Info(msg string, args ...any) {
	fargs := formatArgs(args...)
	l.infoLogger.Printf("%s %s\n", msg, fargs)
}

func (l stdlogger) Error(msg string, args ...any) {
	fargs := formatArgs(args...)
	l.errorLogger.Printf("%s %s\n", msg, fargs)
}

func getArgMap(args ...any) map[string]any {
	result := make(map[string]any)
	for len(args) > 0 {
		key := fmt.Sprint(args[0])
		args = args[1:]

		value := "(missing)"
		if len(args) != 0 {
			value = fmt.Sprint(args[0])
			args = args[1:]

		}
		result[key] = value
	}
	return result
}

func formatArgs(args ...any) string {
	var pairs []string
	for len(args) > 0 {
		key := fmt.Sprint(args[0])
		args = args[1:]

		value := "(missing)"
		if len(args) != 0 {
			value = fmt.Sprint(args[0])
			args = args[1:]

		}
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(pairs, " ")
}
