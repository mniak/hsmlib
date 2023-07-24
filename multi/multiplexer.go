package multi

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/internal/noop"
)

const (
	KConnectionID = "connID"
	KRequestID    = "reqID"
	KError        = "error"
)

type Multiplexer struct {
	listenAddress string
	targetAddress string
	timeout       time.Duration
	logger        hsmlib.Logger
	server        hsmlib.PacketServer

	out           hsmlib.Reactor
	connectionIDs atomic.Uint64
}

type MultiplexerOption func(m *Multiplexer)

func WithTimeout(timeout time.Duration) MultiplexerOption {
	return func(m *Multiplexer) {
		m.timeout = timeout
	}
}

func WithLogger(logger hsmlib.Logger) MultiplexerOption {
	return func(m *Multiplexer) {
		m.logger = logger
		m.server.Logger = logger
	}
}

var DefaultTimout = 10 * time.Second

func NewMultiplexer(listenAddr, targetAddr string, opts ...MultiplexerOption) *Multiplexer {
	result := &Multiplexer{
		listenAddress: listenAddr,
		targetAddress: targetAddr,
		logger:        noop.Logger(),
	}
	for _, opt := range opts {
		opt(result)
	}
	return result
}

func (m *Multiplexer) Run() error {
	outaddr, err := net.ResolveTCPAddr("tcp", m.targetAddress)
	if err != nil {
		return err
	}
	outConn, err := net.DialTCP("tcp", nil, outaddr)
	if err != nil {
		return err
	}
	defer outConn.Close()

	m.out = hsmlib.NewReactorFromReadWriter(outConn)
	m.out.Start()
	m.logger.Info("Multiplexer started")

	return m.server.ListenAndServe(m.listenAddress, hsmlib.PacketHandlerFunc(m.HandleConnection))
}

func (m *Multiplexer) HandleConnection(ps hsmlib.PacketSender, packet hsmlib.Packet) error {
	connID := m.connectionIDs.Add(1)

	timeoutCtx, cancelCtx := context.WithTimeout(context.Background(), m.timeout)
	defer cancelCtx()

	reply, err := m.out.Post(timeoutCtx, packet.Payload)
	if err != nil {
		if m.logger != nil {
			m.logger.Error("failed to forward message",
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

func (m *Multiplexer) Shutdown() {
	m.server.Shutdown()
}
