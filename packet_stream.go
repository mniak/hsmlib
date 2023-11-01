package hsmlib

import (
	"io"
)

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

func (ps _PacketStream) WithLogger(logger Logger) _PacketStreamWithLogs {
	return _PacketStreamWithLogs{
		inner:  ps,
		Logger: logger,
	}
}

type _PacketStreamWithLogs struct {
	Logger Logger
	inner  _PacketStream
}

func (ps _PacketStreamWithLogs) ReceivePacket() (Packet, error) {
	p, err := ps.inner.ReceivePacket()
	if ps.Logger != nil {
		ps.Logger.Info("Packet received",
			"header", p.Header,
			"payload", p.Payload,
		)
	}
	return p, err
}

func (ps _PacketStreamWithLogs) SendPacket(p Packet) error {
	err := ps.inner.SendPacket(p)
	if ps.Logger != nil {
		ps.Logger.Info("Packet sent",
			"header", p.Header,
			"payload", p.Payload,
		)
	}
	return err
}
