package handlers

import (
	"github.com/mniak/hsmlib"
)

func NoOpDecorator(logger hsmlib.Logger) hsmlib.CommandHandlerDecoratorFunc {
	return func(inner hsmlib.CommandHandler, cmd hsmlib.CommandWithHeader) (hsmlib.Response, error) {
		return inner.Handle(cmd)
	}
}
