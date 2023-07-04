package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mniak/hsmlib"
	"github.com/spf13/cobra"
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
				Header: binary.BigEndian.AppendUint32(nil, gofakeit.Uint32()),
				Payload: hsmlib.RawCommand{
					Code: "B2",
					Data: []byte(message),
				}.Bytes(),
			}
			err = hsmlib.SendPacket(conn, packet)
			if err != nil {
				log.Fatalln(err)
			}

			reply, err := hsmlib.ReceiveCommand(conn)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("Header: % 2X\n", reply.Header)
			fmt.Printf("Response Code: %s\n", reply.Code)
			fmt.Printf("Data: %s\n", string(reply.Data))

			if !bytes.Equal(reply.Header, packet.Header) {
				log.Fatalln("reply header invalid")
			}

			if len(reply.Data) != len([]byte(message))+2 {
				log.Fatalln("reply data size invalid")
			}

			if !bytes.Equal(reply.Data[2:], []byte(message)) {
				log.Fatalln("reply data message invalid")
			}
		},
	}
	mainCmd.Execute()
}
