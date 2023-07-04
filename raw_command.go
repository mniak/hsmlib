package hsmlib

import (
	"bytes"
	"errors"
	"io"
)

type RawCommand struct {
	Code string
	Data []byte
}

func (c RawCommand) WithHeader(header []byte) Command {
	return Command{
		RawCommand: c,
		Header:     header,
	}
}

func CalculateResponseCode(commandCode string) string {
	if len(commandCode) < 2 {
		return commandCode
	}
	b := []byte(commandCode)
	b[1]++
	return string(commandCode)
}

func (c RawCommand) Bytes() []byte {
	return append([]byte(c.Code), c.Data...)
}

type Response struct {
	ErrorCode string
	Data      []byte
}

func (r Response) serializeResponse(originalCommandCode string) []byte {
	responseCode := CalculateResponseCode(originalCommandCode)
	var buf bytes.Buffer
	buf.WriteString(responseCode)
	buf.WriteString(r.ErrorCode)
	buf.Write(r.Data)
	return buf.Bytes()
}

func ParseRawCommand(b []byte) (RawCommand, error) {
	const codeLength = 2
	if len(b) < codeLength {
		return RawCommand{}, errors.New("command data is shorter than the length of a command code")
	}

	cmd := RawCommand{
		Code: string(b[:codeLength]),
		Data: b[codeLength:],
	}
	return cmd, nil
}

func ReceiveRawCommand(r io.Reader) (RawCommand, error) {
	packet, err := ReceivePacket(r)
	if err != nil {
		return RawCommand{}, err
	}
	cmd, err := ParseRawCommand(packet.Payload)
	return cmd, err
}
