package translate

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/mniak/hsmlib/cmd/send/internal/app"
	"github.com/mniak/hsmlib/errcode"
	"github.com/mniak/krypton"
	"github.com/mniak/krypton/encoding/hex"
	"github.com/mniak/krypton/futurex"
	"github.com/mniak/krypton/thales"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

type DescriptiveBuilder struct {
	msg bytes.Buffer
	log bytes.Buffer
}

func (b *DescriptiveBuilder) WriteComment(comment string) {
	fmt.Fprintln(&b.log)
	fmt.Fprintf(&b.log, "//// %s\n", comment)
}

func (b *DescriptiveBuilder) Write(data, description string) {
	b.msg.WriteString(data)
	fmt.Fprintf(&b.log, "%s\t// %s\n", data, description)
}

func (b *DescriptiveBuilder) Bytes() []byte {
	return b.msg.Bytes()
}

func (b *DescriptiveBuilder) String() string {
	return b.log.String()
}

func (b *DescriptiveBuilder) Reset() {
	b.msg.Reset()
	b.log.Reset()
}

func cmdTranslateKey() cobra.Command {
	cmd := cobra.Command{
		Use:     "key <Key> <Type> <Target Scheme> <Target Usage>",
		Aliases: []string{"BW"},
		Short:   "Translate key from one LMK to another",
		Args:    cobra.ExactArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			const exportability = "S" // S: Sensitive
			const modeOfUse = "N"     // N: No restrictions
			const keyVersion = "00"   // 00: No versioning
			const encryptdKEK = "BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB"

			sourceKeyScheme, err := thales.ParseKeyScheme(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			clearKey := hex.MustParseString(args[1])
			keyType, err := krypton.KeyTypeByName(args[2])
			if err != nil {
				log.Fatalln(err)
			}
			targetKeyScheme := args[3]
			keyUsage := args[4]

			const keySeparation = false
			keyTypeCode := keyType.Code(keySeparation)
			modifier := lo.Must(futurex.ModifierFromKeyCode(keyTypeCode, keySeparation))

			var hb strings.Builder
			fmt.Fprintf(&hb, "//// Key Info\n")
			fmt.Fprintf(&hb, "Key:\n")
			kcv := lo.Must(krypton.ComputeKCV(clearKey))
			fmt.Fprintf(&hb, "    Cleartext: %2X (KCV=%2X)\n", clearKey, kcv)
			krypton.SetParityOdd(clearKey)
			fmt.Fprintf(&hb, "    With odd parity: %2X\n", clearKey)
			fmt.Fprintf(&hb, "Futurex Cryptogram\n")
			fmt.Fprintf(&hb, "    Master Key: Triple Length MFK Test Key\n")
			mfk := futurex.MFKTriple()
			fmt.Fprintf(&hb, "        Value: %2X\n", mfk)
			fmt.Fprintf(&hb, "    Modifier: %X\n", modifier)
			encryptedCVK := lo.Must(futurex.EncryptKey(mfk, modifier, clearKey))
			fmt.Fprintf(&hb, "    Cryptogram: %2X\n", encryptedCVK)
			fmt.Fprintf(&hb, "Target Key Block\n")
			fmt.Fprintf(&hb, "    Key Scheme: %s\n", targetKeyScheme)
			fmt.Fprintf(&hb, "    Type: %s\n", keyType.Name())
			fmt.Fprintf(&hb, "        Variant Key Type: %s\n", keyTypeCode)
			fmt.Fprintf(&hb, "        Key Usage: %s\n", keyUsage)

			var b DescriptiveBuilder
			b.WriteComment("Command BW: Translate Keys between Master File Keys")
			b.Write("BW", "Command Indicator")
			b.Write("FF", "Key Code - FF to use Key Type defined after delimiter")
			b.Write("1", "Key Length Flag - 1: Double Length")
			b.Write(string(sourceKeyScheme)+hex.ToString(encryptedCVK), "Encrypted Key under the Old MFK")

			b.WriteComment("Conditional 3-digit type code section")
			b.Write(";", "Delimiter")
			b.Write(keyTypeCode, "Key Type - "+keyType.Name())

			b.WriteComment("Target key scheme section")
			b.Write(";", "Delimiter")
			b.Write("0", "Reserved field, must be 0")
			b.Write(targetKeyScheme, "Key Scheme for encrypting key under MFK")
			b.Write("0", "Reserved field, must be 0. Key Check Value Calculation Method")

			// b.WriteComment("Conditional Major Key Identifier section")
			// b.Write("%", "Delimiter")
			// b.Write("00", "Major Key Identifier")

			if targetKeyScheme == "R" || targetKeyScheme == "S" {
				b.WriteComment("Section required if generating Key Block")
				b.Write("#", "Delimiter")
				b.Write(keyUsage, "Key usage")
				b.Write(modeOfUse, "Mode of use")
				b.Write(keyVersion, "Key version")
				b.Write(exportability, "Exportability")

				b.WriteComment("Section of the Optional Blocks")
				b.Write("00", "Number of Optional Blocks")
			}

			b.WriteComment("KCV Section")
			b.Write("!", "Delimiter")
			b.Write("1", "Key Check Return Flag - If present must be '1'")
			b.Write("1", "Key Check Value Type - '1': 6 digit KCV")

			app.Logger().Info("Running command",
				"header", hb.String(),
				"command", string(b.Bytes()),
				"description", b.String(),
			)

			reply := app.SendPacketPayload(b.Bytes())
			if reply.ErrorCode() != errcode.NoError {
				app.Logger().Error("Invalid error code",
					"error_code", reply.ErrorCode(),
				)
			}
			// keyBlock := string(reply.Data())

			// // ------------------- EXPORT (A8) ----------------------
			// b.WriteComment("EXPORT KEY - A8")
			// b.Write("A8", "A8 - Export Key")
			// b.Write("FFF", "Key type - FFF Reserved for Key Blocks")
			// b.Write(encryptdKEK, "KEK under MFK")
			// b.Write(keyBlock, "Working key encrypted under MFK")
			// b.Write("R", "Export Key scheme")

			// b.WriteComment("Export Key Block under MFK to Key Block under KEK")
			// b.Write("&", "Delimiter")
			// b.Write("N", "Exportability")
			// b.Write("!", "Delimiter")
			// b.Write("C", "Export key block Version ID - C: 3DES Variant")

			// app.Logger().Info("Running command",
			// 	"header", hb.String(),
			// 	"command", string(b.Bytes()),
			// 	"description", b.String(),
			// )

			// reply = app.SendPacketPayload(b.Bytes())
			// if reply.ErrorCode() != errcode.NoError {
			// 	app.Logger().Error("Invalid error code",
			// 		"error_code", reply.ErrorCode(),
			// 	)
			// }
		},
	}

	return cmd
}
