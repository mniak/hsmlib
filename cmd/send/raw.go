package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/cmd/send/internal/app"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func rawCommand() *cobra.Command {
	var flagAutoHeader bool
	var flagHex bool
	cmd := cobra.Command{
		Use:   "raw [--auto-header] <DATA>",
		Short: "Sends raw data to an HSM",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			data := []byte(args[0])
			if flagHex {
				var err error
				data, err = hex.DecodeString(string(data))
				if err != nil {
					log.Fatalln("Invalid hexadecimal data")
				}
			}

			var reply hsmlib.ResponseWithHeader
			if flagAutoHeader {
				packet := hsmlib.Packet{
					Header:  binary.BigEndian.AppendUint32(nil, gofakeit.Uint32()),
					Payload: data,
				}
				slog.Info("Sending raw packet:",
					"header", fmt.Sprintf("%2X", packet.Header),
					"payload", fmt.Sprintf("%2X", packet.Payload),
				)
				reply = app.SendPacket(packet)
			} else {
				slog.Info("Sending raw frame:",
					"data", fmt.Sprintf("%2X", data),
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
