package main

import (
	"errors"
	"io"
)

type SimpleIDWrapper struct {
	IDLength int
	Inner    Protocol[[]byte]
}

type PacketWithID struct {
	ID   []byte
	Data []byte
}

func (p PacketWithID) Bytes() []byte {
	return append(p.ID, p.Data...)
}

func (m SimpleIDWrapper) Send(w io.Writer, packet PacketWithID) error {
	bytes := packet.Bytes()
	err := m.Inner.Send(w, bytes)
	return err
}

func (p SimpleIDWrapper) Receive(r io.Reader) (PacketWithID, error) {
	bytes, err := p.Inner.Receive(r)
	if err != nil {
		return PacketWithID{}, err
	}
	if len(bytes) < p.IDLength {
		return PacketWithID{}, errors.New("packet is shorter than the length of an ID")
	}
	packet := PacketWithID{
		ID:   bytes[:p.IDLength],
		Data: bytes[p.IDLength:],
	}
	return packet, nil
}
