package misc

import (
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func RegisterCommands(parent *cobra.Command) {
	parent.AddCommand(lo.ToPtr(cmdEcho()))
	parent.AddCommand(lo.ToPtr(cmdKCV()))
}
