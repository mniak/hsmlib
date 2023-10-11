package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/cmd/send/internal/app"
	"github.com/mniak/hsmlib/commands"
	"github.com/spf13/cobra"
)

func echoCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "echo <MESSAGE>",
		Aliases: []string{"B2"},
		Short:   "Sends an echo command to and HSM",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			message := args[0]

			packet := hsmlib.Packet{
				Header:  makeHeader(),
				Payload: hsmlib.CommandBytes(commands.MakeEcho(message)),
			}
			reply := app.SendPacket(packet)

			var failed bool
			if reply.ErrorCode != "00" {
				fmt.Fprintf(os.Stderr, "Echo received invalid error code: %q\n", reply.ErrorCode)
				failed = true
			}
			if !bytes.Equal(reply.Data, []byte(message)) {
				fmt.Fprintf(os.Stderr, "Echo received invalid message: %q\n", string(reply.Data))
				failed = true
			}
			if failed {
				os.Exit(1)
			} else {
				fmt.Printf("Echo success!")
			}
		},
	}

	return &cmd
}
