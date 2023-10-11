package hsmlib

import (
	"errors"
	"io"
)

type Command interface {
	Code() []byte
	Data() []byte
}

type RawCommand struct {
	RawCode string
	RawData []byte
}

func (cmd RawCommand) Code() []byte {
	return []byte(cmd.RawCode)
}

func (cmd RawCommand) Data() []byte {
	return cmd.RawData
}

func (c RawCommand) WithHeader(header []byte) CommandWithHeader {
	return CommandWithHeader{
		Command: c,
		Header:  header,
	}
}

func CommandBytes(cmd Command) []byte {
	return append(cmd.Code(), cmd.Data()...)
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

func ParseCommand(b []byte) (RawCommand, error) {
	const codeLength = 2
	if len(b) < codeLength {
		return RawCommand{}, errors.New("command data is shorter than the length of a command code")
	}

	cmd := RawCommand{
		RawCode: string(b[:codeLength]),
		RawData: b[codeLength:],
	}
	return cmd, nil
}

func ReceiveRawCommand(r io.Reader) (Command, error) {
	packet, err := ReceivePacket(r)
	if err != nil {
		return RawCommand{}, err
	}
	cmd, err := ParseCommand(packet.Payload)
	return cmd, err
}
