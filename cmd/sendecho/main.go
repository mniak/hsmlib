package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
	"net"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/commands"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func main() {
	mainCmd := cobra.Command{
		Use:   "sendecho <TARGET> <MESSAGE>",
		Short: "Sends an echo command to and HSM",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			target := args[0]
			message := args[1]

			conn, err := net.Dial("tcp", target)
			if err != nil {
				log.Fatalln(err)
			}
			defer conn.Close()

			packet := hsmlib.Packet{
				Header:  binary.BigEndian.AppendUint32(nil, gofakeit.Uint32()),
				Payload: commands.Echo(message).Bytes(),
			}
			err = hsmlib.SendPacket(conn, packet)
			if err != nil {
				log.Fatalln(err)
			}

			reply, err := hsmlib.ReceiveResponse(conn)
			if err != nil {
				log.Fatalln(err)
			}

			slog.Info("Received response",
				"header", hex.EncodeToString(reply.Header),
				"response_code", reply.ResponseCode,
				"error_code", reply.ErrorCode,
				"data", string(reply.Data),
			)

			if !bytes.Equal(reply.Header, packet.Header) {
				log.Fatalln("reply header invalid")
			}

			if len(reply.Data) != len([]byte(message)) {
				log.Fatalln("reply data size invalid")
			}

			if !bytes.Equal(reply.Data, []byte(message)) {
				log.Fatalln("reply data message invalid")
			}
		},
	}
	mainCmd.Execute()
}
