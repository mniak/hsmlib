package hsmlib

import "strings"

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

func (r *CommandRouter) AddFallbackHandler(h CommandHandler) *CommandRouter {
	r.fallbackHandler = h
	return r
}

func (r *CommandRouter) AddHandler(commandCode string, h CommandHandler) *CommandRouter {
	r.handlers[strings.ToUpper(commandCode)] = h
	return r
}

func (r *CommandRouter) FindHandler(req CommandWithHeader) CommandHandler {
	if handler, found := r.handlers[strings.ToUpper(req.Code)]; found {
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
