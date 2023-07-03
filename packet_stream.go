package hsmlib

import "io"

type PacketSender interface {
	SendPacket(Packet) error
}
type PacketReceiver interface {
	ReceivePacket() (Packet, error)
}
type PacketStream interface {
	PacketSender
	PacketReceiver
}

func NewPacketStream(rw io.ReadWriter) _PacketStream {
	return _PacketStream{
		rw: rw,
	}
}

type _PacketStream struct {
	rw io.ReadWriter
}

func (ps _PacketStream) ReceivePacket() (Packet, error) {
	return ReceivePacket(ps.rw)
}

func (ps _PacketStream) SendPacket(p Packet) error {
	return SendPacket(ps.rw, p)
}
