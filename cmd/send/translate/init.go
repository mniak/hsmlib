package translate

import (
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func RegisterCommands(parent *cobra.Command) {
	translate := cobra.Command{
		Use:   "translate",
		Short: "Translate PINs and keys",
	}
	translate.AddCommand(lo.ToPtr(cmdTranslateKey()))
	parent.AddCommand(&translate)
}
