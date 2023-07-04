package hsmlib

import (
	"io"
)

type Command struct {
	Header []byte
	RawCommand
}

func ReceiveCommand(r io.Reader) (Command, error) {
	packet, err := ReceivePacket(r)
	if err != nil {
		return Command{}, err
	}
	rawCmd, err := ParseRawCommand(packet.Payload)
	cmd := Command{
		RawCommand: rawCmd,
		Header:     packet.Header,
	}
	return cmd, err
}
