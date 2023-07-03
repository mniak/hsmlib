package main

import (
	"context"
	"io"
	"log"
)

type Reactor struct {
	IDManager IDManager
	Protocol  Protocol[PacketWithID]
	RW        io.ReadWriter
}

func (m *Reactor) Start() {
	go func() {
		for {
			packet, err := m.Protocol.Receive(m.RW)
			if err != nil {
				log.Println("error handling packet", err)
				continue
			}

			channel, found := m.IDManager.FindChannel(packet.ID)
			if !found {
				log.Printf("callback channel not found for id %02X\n", packet.ID)
				continue
			}

			go func() {
				channel <- packet.Data
			}()
		}
	}()
}

func (m *Reactor) Post(ctx context.Context, data []byte) ([]byte, error) {
	id, ch := m.IDManager.NewID()
	err := m.Protocol.Send(m.RW, PacketWithID{
		ID:   id,
		Data: data,
	})
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

func NewReactor(rw io.ReadWriter) Reactor {
	idWrap := NewHSMProtocol()
	reactor := Reactor{
		IDManager: &SequentialUint32IDManager{},
		Protocol:  idWrap,
		RW:        rw,
	}
	return reactor
}
