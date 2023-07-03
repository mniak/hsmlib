package multi

import (
	"sync/atomic"
	"time"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/internal/noop"
	"go.uber.org/multierr"
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

	out           Reactor
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
	}
}

var DefaultTimout = 10 * time.Second

func NewMultiplexer(listenAddr, targetAddr string, opts ...MultiplexerOption) *Multiplexer {
	m := &Multiplexer{
		listenAddress: listenAddr,
		targetAddress: targetAddr,
	}
	for _, opt := range opts {
		opt(m)
	}

	if m.logger == nil {
		m.logger = noop.Logger()
	}
	m.server = hsmlib.NewPacketServer(m.logger)
	return m
}

func (m *Multiplexer) Run() error {
	reactor := NewResilientReactor(m.targetAddress, m.logger)
	if err := reactor.Start(); err != nil {
		return err
	}
	m.out = reactor

	m.logger.Info("Multiplexer started")

	packetHandler := hsmlib.PacketHandler(hsmlib.PacketHandlerFunc(m.HandleConnection))
	return hsmlib.ListenAndServeI(m.server, m.listenAddress, packetHandler)
}

func (m *Multiplexer) HandleConnection(in hsmlib.PacketSender, packet hsmlib.Packet) error {
	connID := m.connectionIDs.Add(1)

	reply, err := m.out.Post(packet.Payload)
	if err != nil {
		m.logger.Error("failed to forward message",
			KConnectionID, connID,
			KRequestID, packet.Header,
			KError, err,
		)
		err = m.sendFailureResponse(in, packet, err)
		return err
	}
	err = in.SendPacket(hsmlib.Packet{
		Header:  packet.Header,
		Payload: reply,
	})
	return err
}

func (m *Multiplexer) sendFailureResponse(in hsmlib.PacketSender, packet hsmlib.Packet, err error) error {
	cmd, err2 := hsmlib.ParseCommand(packet.Payload)
	if err2 != nil {
		return multierr.Combine(err, err2)
	}
	resp := hsmlib.Response{
		ErrorCode: "99",
		Data:      []byte(err.Error()),
	}
	respPacket := resp.WithCode(cmd.Code).
		WithHeader(packet.Header).
		AsPacket()

	err2 = in.SendPacket(hsmlib.Packet{
		Header:  packet.Header,
		Payload: respPacket.Bytes(),
	})
	if err2 != nil {
		return multierr.Combine(err, err2)
	}
	return nil
}

func (m *Multiplexer) Shutdown() {
	m.server.Shutdown()
}
