package main

import (
	"fmt"
	"time"

	"github.com/mniak/hsmlib"
	"golang.org/x/exp/slog"
)

func RunHSMMock(address string) error {
	return hsmlib.ListenAndServe(address, hsmlib.CommandHandlerFunc(func(cmd hsmlib.CommandWithHeader) (reply hsmlib.Response, err error) {
		slog.Info("Command received",
			"command.header", fmt.Sprintf("[% 2X]", cmd.Header),
			"command.code", cmd.Code,
			"command.data", cmd.Data,
		)
		time.Sleep(time.Second * 1)

		resp := hsmlib.Response{
			ErrorCode: "00",
			Data:      cmd.Data,
		}
		return resp, nil
	}))
}
