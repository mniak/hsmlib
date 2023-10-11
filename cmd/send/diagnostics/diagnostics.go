package diagnostics

import (
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func RegisterCommands(cmd *cobra.Command) {
	cmd.AddCommand(lo.ToPtr(cmdHealthcheck()))
}
