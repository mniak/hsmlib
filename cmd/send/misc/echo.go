package misc

import (
	"bytes"
	"fmt"
	"os"

	"github.com/mniak/hsmlib/cmd/send/internal/app"
	"github.com/mniak/hsmlib/cmds"
	"github.com/mniak/hsmlib/errcode"
	"github.com/spf13/cobra"
)

func cmdEcho() cobra.Command {
	cmd := cobra.Command{
		Use:     "echo <MESSAGE>",
		Aliases: []string{"B2"},
		Short:   "Sends an echo command to an HSM",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			message := args[0]
			reply := app.SendCommand(cmds.MakeEcho(message))

			var failed bool
			if reply.ErrorCode() != errcode.NoError {
				app.Logger().Error("Invalid error code",
					"error_code", fmt.Sprintf("%q", reply.ErrorCode()),
				)
				failed = true
			}
			if !bytes.Equal(reply.Data(), []byte(message)) {
				app.Logger().Error("Invalid message",
					"message", reply.Data(),
				)
				failed = true
			}
			if failed {
				os.Exit(1)
			} else {
				app.Logger().Info("Success!")
			}
		},
	}

	return cmd
}
