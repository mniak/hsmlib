package main

import (
	"github.com/spf13/cobra"
)

func translateCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "translate <SUBCOMMAND>",
		Short: "Translate keys or values",
	}

	return &cmd
}
