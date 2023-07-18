package main

import (
	"context"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/mniak/hsmlib"
	"github.com/samber/lo"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

const (
	KConnectionID = "connID"
	KRequestID    = "reqID"
	KError        = "error"
)

type Multiplexer struct {
	ListenAddress string
	TargetAddress string
	Timeout       time.Duration
	Logger        Logger

	out           hsmlib.Reactor
	connectionIDs atomic.Uint64
}

func (m *Multiplexer) Run() error {
	outaddr := lo.Must(net.ResolveTCPAddr("tcp", m.TargetAddress))
	outConn := lo.Must(net.DialTCP("tcp", nil, outaddr))
	defer outConn.Close()

	m.out = hsmlib.NewReactorFromReadWriter(outConn)
	m.out.Start()
	log.Println("reactor started")

	return hsmlib.ListenAndServeRaw(m.ListenAddress, hsmlib.PacketHandlerFunc(m.HandleConnection))
}

func (m *Multiplexer) HandleConnection(ps hsmlib.PacketSender, packet hsmlib.Packet) error {
	connID := m.connectionIDs.Add(1)

	timeoutCtx, cancelCtx := context.WithTimeout(context.Background(), m.Timeout)
	defer cancelCtx()

	reply, err := m.out.Post(timeoutCtx, packet.Payload)
	if err != nil {
		if m.Logger != nil {
			m.Logger.Error("failed to forward message",
				KConnectionID, connID,
				KRequestID, packet.Header,
				KError, err,
			)
		}
		return err
	}
	err = ps.SendPacket(hsmlib.Packet{
		Header:  packet.Header,
		Payload: reply,
	})
	return err
}
