package hsmlib

import (
	"context"
	"io"
	"log"
	"net"
)

type Reactor struct {
	IDManager IDManager
	Stream    io.ReadWriter
	done      chan struct{}
}

func NewReactorFromReadWriter(rw io.ReadWriter) Reactor {
	reactor := Reactor{
		IDManager: &SequentialIDManager{},
		Stream:    rw,
	}
	return reactor
}

func NewReactor(target string) (Reactor, error) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return Reactor{}, err
	}
	return NewReactorFromReadWriter(conn), nil
}

func (m *Reactor) Start() {
	m.done = make(chan struct{})
	go func() {
		for {
			select {
			case <-m.done:
				return
			default:
				m.receiveOnePacket()
			}
		}
	}()
}

func (m *Reactor) receiveOnePacket() {
	packet, err := ReceivePacket(m.Stream)
	if err != nil {
		log.Println("error receiving frame", err)
		return
	}
	channel, found := m.IDManager.FindChannel(packet.Header)
	if !found {
		log.Printf("callback channel not found for id %02X\n", packet.Header)
		return
	}

	go func() {
		channel <- packet.Payload
		close(channel)
	}()
}

func (m *Reactor) Stop() {
	close(m.done)
}

func (m *Reactor) Post(ctx context.Context, data []byte) ([]byte, error) {
	id, ch := m.IDManager.NewID()
	packet := Packet{
		Header:  id,
		Payload: data,
	}
	err := SendPacket(m.Stream, packet)
	if err != nil {
		return nil, err
	}

	select {
	case response := <-ch:
		return response, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
