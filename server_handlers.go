package hsmlib

type CommandHandler interface {
	Handle(cmd CommandWithHeader) (Response, error)
}

type CommandHandlerFunc func(cmd CommandWithHeader) (Response, error)

func (h CommandHandlerFunc) Handle(cmd CommandWithHeader) (Response, error) {
	return h(cmd)
}

type CommandHandlerDecorator interface {
	Handle(inner CommandHandler, cmd CommandWithHeader) (Response, error)
}

type CommandHandlerDecoratorFunc func(inner CommandHandler, cmd CommandWithHeader) (Response, error)

func (hd CommandHandlerDecoratorFunc) Handle(inner CommandHandler, cmd CommandWithHeader) (Response, error) {
	return hd(inner, cmd)
}
