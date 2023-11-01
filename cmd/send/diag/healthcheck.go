package diag

import (
	"fmt"

	"github.com/mniak/hsmlib/cmd/send/internal/app"
	"github.com/mniak/hsmlib/cmds"
	"github.com/mniak/hsmlib/errcode"
	"github.com/spf13/cobra"
)

func cmdHealthcheck() cobra.Command {
	cmd := cobra.Command{
		Use:     "healthcheck",
		Aliases: []string{"JK", "health"},
		Short:   "Get instantaneous health check status",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			reply := app.SendCommand(cmds.MakeHealthcheck())
			if reply.ErrorCode() != errcode.NoError {
				app.Logger().Error("Invalid error code",
					"error_code", fmt.Sprintf("%q", reply.ErrorCode()),
				)
			}
		},
	}

	return cmd
}
