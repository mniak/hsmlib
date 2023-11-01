package main

import (
	"encoding/hex"
	"log"
	"strings"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/cmd/send/internal/app"
	"github.com/spf13/cobra"
)

func cmdRaw() *cobra.Command {
	var flagAutoHeader bool
	var flagHex bool
	cmd := cobra.Command{
		Use:   "raw [--auto-header] <DATA>",
		Short: "Sends raw data to an HSM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			data := []byte(strings.Join(args, ""))
			if flagHex {
				var err error
				data, err = hex.DecodeString(string(data))
				if err != nil {
					log.Fatalln("Invalid hexadecimal data")
				}
			}

			var reply hsmlib.ResponseWithHeader
			if flagAutoHeader {
				app.Logger().Info("Sending raw packet payload:",
					"data", data,
				)
				reply = app.SendPacketPayload(data)
			} else {
				app.Logger().Info("Sending raw frame:",
					"data", data,
				)
				reply = app.SendFrame(data)
			}
			_ = reply
		},
	}
	cmd.Flags().BoolVar(&flagAutoHeader, "auto-header", false, "Generate and append a 4 byte header")
	cmd.Flags().BoolVar(&flagHex, "hex", false, "Consider the data to be in hexadecimal")
	return &cmd
}
