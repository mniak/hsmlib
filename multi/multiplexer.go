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
	Logger  hsmlib.Logger
	Timeout time.Duration

	server        hsmlib.PacketServer
	out           Reactor
	connectionIDs atomic.Uint64
}

var DefaultTimeout = 10 * time.Second

func (m *Multiplexer) ensureInit() {
	if m.Logger == nil {
		m.Logger = noop.Logger()
	}
	m.server = hsmlib.NewPacketServer(m.Logger)
}

func (m *Multiplexer) Run(listenAddr string, dialer Dialer) error {
	m.ensureInit()

	reactor := ResilientReactor{
		Logger: m.Logger,
	}
	if err := reactor.Start(dialer); err != nil {
		return err
	}
	m.out = &reactor

	m.Logger.Info("Multiplexer started")

	packetHandler := hsmlib.PacketHandler(hsmlib.PacketHandlerFunc(m.HandleConnection))
	return hsmlib.ListenAndServeI[hsmlib.PacketHandler](m.server, listenAddr, packetHandler)
}

func (m *Multiplexer) HandleConnection(in hsmlib.PacketSender, packet hsmlib.Packet) error {
	connID := m.connectionIDs.Add(1)

	reply, err := m.out.Post(packet.Payload)
	if err != nil {
		m.Logger.Error("failed to forward message",
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
	respPacket := resp.ForCommandCode(cmd.Code()).
		WithHeader(packet.Header).
		AsPacket()

	err2 = in.SendPacket(respPacket)
	if err2 != nil {
		return multierr.Combine(err, err2)
	}
	return nil
}

func (m *Multiplexer) Shutdown() {
	m.server.Shutdown()
}
