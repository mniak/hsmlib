package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func main() {
	var (
		flagHost string
		flagPort uint16
	)
	mainCmd := cobra.Command{
		Use: "mock [--host <HOST>] [--port <PORT>]",
		Run: func(cmd *cobra.Command, args []string) {
			address := fmt.Sprintf("%s:%d", flagHost, flagPort)
			err := RunHSMMock(address)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
	mainCmd.PersistentFlags().BoolP("help", "", false, "Help for this command")
	mainCmd.Flags().StringVarP(&flagHost, "host", "h", "127.0.0.1", "Hostname for listening (don't include port number)")
	mainCmd.Flags().Uint16VarP(&flagPort, "port", "p", 1501, "Port for listening")

	mainCmd.Execute()
}
