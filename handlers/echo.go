package handlers

import (
	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/cmds"
	"github.com/mniak/hsmlib/errcode"
	"github.com/mniak/hsmlib/internal/noop"
)

type EchoHandler struct {
	Logger hsmlib.Logger
}

func (h *EchoHandler) guard() {
	if h.Logger == nil {
		h.Logger = noop.Logger()
	}
}

func (h *EchoHandler) Handle(cmd hsmlib.CommandWithHeader) (hsmlib.Response, error) {
	h.guard()

	echo, err := cmds.EchoFromCommand(cmd)
	if err != nil {
		return nil, err
	}

	h.Logger.Info("Echo command received",
		"message", echo.Message,
	)

	return cmds.MakeEchoResponse(errcode.NoError, echo.Message), nil
}
