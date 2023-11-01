package hsmlib

import (
	"strings"
)

type CommandRouter struct {
	handlers        map[string]CommandHandler
	fallbackHandler CommandHandler
}

var DefaultFallbackCommandHandler CommandHandlerFunc = func(cmd CommandWithHeader) (Response, error) {
	return CommandDisabledResponse()
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		handlers: make(map[string]CommandHandler),
	}
}

func (r *CommandRouter) AddFallbackHandler(handler CommandHandler) *CommandRouter {
	r.fallbackHandler = handler
	return r
}

func (r *CommandRouter) AddHandler(commandCode string, handler CommandHandler) *CommandRouter {
	r.handlers[strings.ToUpper(commandCode)] = handler
	return r
}

func (r *CommandRouter) DecorateHandler(commandCode string, decorator CommandHandlerDecorator) *CommandRouter {
	commandCode = strings.ToUpper(commandCode)
	inner := r.handlers[commandCode]
	r.handlers[commandCode] = CommandHandlerFunc(func(cmd CommandWithHeader) (Response, error) {
		return decorator.Handle(inner, cmd)
	})
	return r
}

func (r *CommandRouter) DecorateHandlerFn(commandCode string, decoratorFn CommandHandlerDecoratorFunc) *CommandRouter {
	return r.DecorateHandler(commandCode, decoratorFn)
}

func (r *CommandRouter) FindHandler(req CommandWithHeader) CommandHandler {
	cmdString := strings.ToUpper(string(req.Code()))
	if handler, found := r.handlers[cmdString]; found {
		return handler
	}
	if r.fallbackHandler != nil {
		return r.fallbackHandler
	}
	return DefaultFallbackCommandHandler
}

func (r *CommandRouter) Handle(cmd CommandWithHeader) (Response, error) {
	h := r.FindHandler(cmd)
	return h.Handle(cmd)
}
