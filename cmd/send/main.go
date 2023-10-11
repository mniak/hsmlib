package main

import (
	"encoding/binary"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/hsmlib/cmd/send/diagnostics"
	"github.com/mniak/hsmlib/cmd/send/internal/app"
	"github.com/spf13/cobra"
)

func main() {
	var flagTarget string
	var flagTLS bool
	var flagClientCertFile string
	var flagClientKeyFile string
	var flagSkipVerify bool

	defer app.Finish()
	cmd := cobra.Command{
		Use:   "send --target <TARGET> [<connection flags>]",
		Short: "Sends an echo command to and HSM",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			app.Connect(flagTarget, flagTLS, flagClientCertFile, flagClientKeyFile, flagSkipVerify)
		},
	}
	cmd.PersistentFlags().StringVar(&flagTarget, "target", "", "Specify the connection target")
	cmd.MarkPersistentFlagRequired("target")
	cmd.PersistentFlags().BoolVar(&flagTLS, "tls", false, "Enable TLS in the connection")
	cmd.PersistentFlags().StringVar(&flagClientCertFile, "client-cert-file", "", "Specify a TLS client certificate file")
	cmd.PersistentFlags().StringVar(&flagClientKeyFile, "client-key-file", "", "Specify a TLS client key file")
	cmd.PersistentFlags().BoolVar(&flagSkipVerify, "skip-verify", false, "Don't verify the target's certificate")

	cmd.AddCommand(rawCommand())
	cmd.AddCommand(echoCommand())

	diagnostics.RegisterCommands(&cmd)

	cmd.Execute()
}

func makeHeader() []byte {
	return binary.BigEndian.AppendUint32(nil, gofakeit.Uint32())
}
