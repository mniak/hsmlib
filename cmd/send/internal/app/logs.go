package app

import (
	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/adapters/stdlib"
	"github.com/mniak/hsmlib/internal/noop"
)

var logger hsmlib.Logger = noop.Logger()

func Verbose(verbose bool) {
	if verbose {
		logger = stdlib.NewLogger("[Sender] ")
	} else {
		logger = noop.Logger()
	}
}

func Logger() hsmlib.Logger {
	return logger
}
