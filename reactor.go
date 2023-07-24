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
	Logger    log.Logger
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

func (r *Reactor) receiveOnePacket() {
	packet, err := ReceivePacket(r.Stream)
	if err != nil {
		r.Logger.Println("error receiving frame", err)
		return
	}
	channel, found := r.IDManager.FindChannel(packet.Header)
	if !found {
		r.Logger.Printf("callback channel not found for id %02X\n", packet.Header)
		return
	}

	go func() {
		channel <- packet.Payload
		close(channel)
	}()
}

func (r *Reactor) Stop() {
	close(r.done)
}

func (r *Reactor) Post(ctx context.Context, data []byte) ([]byte, error) {
	id, ch := r.IDManager.NewID()
	packet := Packet{
		Header:  id,
		Payload: data,
	}
	err := SendPacket(r.Stream, packet)
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
