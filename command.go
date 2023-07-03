package hsmlib

import (
	"errors"
	"io"
)

type Command struct {
	Code string
	Data []byte
}

func (c Command) WithHeader(header []byte) CommandWithHeader {
	return CommandWithHeader{
		Command: c,
		Header:  header,
	}
}

func (c Command) Bytes() []byte {
	return append([]byte(c.Code), c.Data...)
}

type CommandWithHeader struct {
	Header []byte
	Command
}

func ReceiveCommand(r io.Reader) (CommandWithHeader, error) {
	packet, err := ReceivePacket(r)
	if err != nil {
		return CommandWithHeader{}, err
	}
	cmd, err := ParseCommand(packet.Payload)
	cmdH := CommandWithHeader{
		Command: cmd,
		Header:  packet.Header,
	}
	return cmdH, err
}

func ParseCommand(b []byte) (Command, error) {
	const codeLength = 2
	if len(b) < codeLength {
		return Command{}, errors.New("command data is shorter than the length of a command code")
	}

	cmd := Command{
		Code: string(b[:codeLength]),
		Data: b[codeLength:],
	}
	return cmd, nil
}

func ReceiveRawCommand(r io.Reader) (Command, error) {
	packet, err := ReceivePacket(r)
	if err != nil {
		return Command{}, err
	}
	cmd, err := ParseCommand(packet.Payload)
	return cmd, err
}
