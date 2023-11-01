package main

import (
	"fmt"
	"time"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/handlers"
	"golang.org/x/exp/slog"
)

func RunHSMMock(address string) error {
	router := hsmlib.NewCommandRouter()

	router.AddHandler("B2", &handlers.EchoHandler{})
	router.DecorateHandlerFn("B2", func(inner hsmlib.CommandHandler, cmd hsmlib.CommandWithHeader) (hsmlib.Response, error) {
		slog.Info("Command received",
			"command.header", fmt.Sprintf("[% 2X]", cmd.Header),
			"command.code", cmd.Code,
			"command.data", cmd.Data,
		)
		time.Sleep(time.Second * 1)
		return inner.Handle(cmd)
	})

	return hsmlib.ListenAndServe(address, router)
}
