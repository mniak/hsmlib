package main

import (
	"fmt"
	"time"

	"github.com/mniak/hsmlib"
	"golang.org/x/exp/slog"
)

func RunHSMMock(address string) error {
	router := hsmlib.NewCommandRouter()
	router.AddHandler("B2", hsmlib.CommandHandlerFunc(func(cmd hsmlib.CommandWithHeader) (reply hsmlib.Response, err error) {
		slog.Info("Command received",
			"command.header", fmt.Sprintf("[% 2X]", cmd.Header),
			"command.code", cmd.Code,
			"command.data", cmd.Data,
		)
		time.Sleep(time.Second * 1)

		message, remaining, err := hsmlib.ParseWithLengthPrefix4H(cmd.Data())
		if err != nil || len(remaining) > 0 {
			return hsmlib.InvalidInputDataResponse()
		}

		resp := hsmlib.Response{
			ErrorCode: "00",
			Data:      message,
		}
		return resp, nil
	}))

	return hsmlib.ListenAndServe(address, router)
}
