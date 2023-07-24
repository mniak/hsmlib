package main

import (
	"log"
	"time"

	"github.com/mniak/hsmlib/adapters/stdlib"
	"github.com/mniak/hsmlib/multi"
	"github.com/spf13/cobra"
)

func main() {
	var (
		flagListenAddress  string
		flagTargetAddress  string
		flagTimeoutSeconds int
	)
	mainCmd := cobra.Command{
		Use: "multiplexer [--listen <address>] [--target <address>] [--timeout <timeout>]",
		Run: func(cmd *cobra.Command, args []string) {
			m := multi.NewMultiplexer(flagListenAddress, flagTargetAddress,
				multi.WithLogger(stdlib.NewLogger("[Multiplexer] ")),
				multi.WithTimeout(time.Duration(flagTimeoutSeconds)*time.Second),
			)
			if err := m.Run(); err != nil {
				log.Fatalln(err)
			}
		},
	}
	mainCmd.Flags().StringVarP(&flagListenAddress, "listen", "l", "0.0.0.0:1500", "Listen address")
	mainCmd.Flags().StringVarP(&flagTargetAddress, "target", "t", "", "Target address")
	mainCmd.MarkFlagRequired("target")
	mainCmd.Flags().IntVar(&flagTimeoutSeconds, "timeout", 10, "Specify the timeout for requests forwarded")

	mainCmd.Execute()
}
