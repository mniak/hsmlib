package diagnostics

import (
	"github.com/spf13/cobra"
)

func cmdHealthcheck() cobra.Command {
	cmd := cobra.Command{
		Use:     "healthcheck <SUBCOMMAND>",
		Aliases: []string{"JK", "health"},
		Short:   "Get instantaneous health check status",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return cmd
}
