package hsmlib

import (
	"errors"
	"io"
)

type Packet struct {
	Header  []byte
	Payload []byte
}

func (p Packet) Bytes() []byte {
	return append(p.Header, p.Payload...)
}

const HeaderLength = 4

func ParsePacket(data []byte) (Packet, error) {
	if len(data) < HeaderLength {
		return Packet{}, errors.New("packet data is shorter than the length of the header")
	}
	packet := Packet{
		Header:  data[:HeaderLength],
		Payload: data[HeaderLength:],
	}
	return packet, nil
}

func ReceivePacket(r io.Reader) (Packet, error) {
	frame, err := ReceiveFrame(r)
	if err != nil {
		return Packet{}, err
	}

	packet, err := ParsePacket(frame)
	return packet, err
}

func SendPacket(w io.Writer, packet Packet) error {
	return SendFrame(w, packet.Bytes())
}
