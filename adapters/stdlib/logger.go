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

const (
	separator = " "
	format    = "%s [%s]\n"
)

func (l stdlogger) Info(msg string, args ...any) {
	fargs := formatArgs(separator, args...)
	if fargs == "" {
		l.infoLogger.Print(msg)
	} else {
		l.infoLogger.Printf(format, msg, fargs)
	}
}

func (l stdlogger) Error(msg string, args ...any) {
	fargs := formatArgs(separator, args...)
	if fargs == "" {
		l.infoLogger.Print(msg)
	} else {
		l.errorLogger.Printf(format, msg, fargs)
	}
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

type argPair struct {
	key   string
	value any
}

func parseArgs(args ...any) []argPair {
	var pairs []argPair
	for len(args) > 0 {
		pair := argPair{
			key: fmt.Sprint(args[0]),
		}
		args = args[1:]

		if len(args) != 0 {
			pair.value = args[0]
			args = args[1:]

		}
		pairs = append(pairs, pair)
	}
	return pairs
}

func formatArgs(separator string, args ...any) string {
	var sb strings.Builder
	for idx, pair := range parseArgs(args...) {
		if idx != 0 {
			sb.WriteString(separator)
		}
		switch pair.value.(type) {
		case []byte:
			fmt.Fprintf(&sb, `%s/hex=%2X`, pair.key, pair.value)
			sb.WriteString(separator)
			fmt.Fprintf(&sb, `%s/str=%q`, pair.key, pair.value)

		default:
			fmt.Fprintf(&sb, `%s=%s`, pair.key, pair.value)
		}
	}
	return sb.String()
}
