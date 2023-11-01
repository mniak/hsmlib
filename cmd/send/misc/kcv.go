package misc

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/cmd/send/internal/app"
	"github.com/mniak/hsmlib/errcode"
	"github.com/mniak/krypton"
	"github.com/mniak/krypton/thales"
	"github.com/spf13/cobra"
)

func cmdKCV() cobra.Command {
	cmd := cobra.Command{
		Use:     "kcv <Key Type> <Key>",
		Aliases: []string{"KCV", "BU"},
		Short:   "Sends an echo command to an HSM",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			keyType, err := krypton.KeyTypeByName(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			keyScheme, err := thales.ParseKeyScheme(args[1][:1])
			if err != nil {
				log.Fatalln(err)
			}
			key := args[1]

			var bb bytes.Buffer
			fmt.Fprintf(&bb, "FF")
			switch keyScheme {
			case thales.KeySchemeU:
				fmt.Fprint(&bb, "1")
			case thales.KeySchemeT:
				fmt.Fprint(&bb, "2")
			case thales.KeyScheme('S'), thales.KeyScheme('R'):
				fmt.Fprint(&bb, "F")
			default:
				log.Fatalln("Unknown length for the specified key scheme")
			}
			fmt.Fprint(&bb, key)
			fmt.Fprint(&bb, ";")
			switch keyScheme {
			case thales.KeySchemeU, thales.KeySchemeT:
				fmt.Fprint(&bb, keyType.Code(false))
			case thales.KeyScheme('S'), thales.KeyScheme('R'):
				fmt.Fprint(&bb, "FFF")
			default:
				log.Fatalln("Unknown key type code for the specified key scheme")
			}
			reply := app.SendCommand(hsmlib.RawCommand{
				RawCode: "BU",
				RawData: bb.Bytes(),
			})

			if reply.ErrorCode() != errcode.NoError {
				app.Logger().Error("Invalid error code",
					"error_code", fmt.Sprintf("%q", reply.ErrorCode()),
				)
				os.Exit(1)
			}

			app.Logger().Info("Recive KCV response",
				"kcv", reply.Data(),
			)
		},
	}

	return cmd
}
