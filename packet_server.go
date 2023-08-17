package hsmlib

import (
	"fmt"
	"io"
	"net"

	"github.com/mniak/hsmlib/internal/noop"
	"github.com/pkg/errors"
)

type PacketServer interface {
	Serve(listener net.Listener, handler PacketHandler) (serveError error)
	Shutdown()
}

type _PacketServer struct {
	logger Logger
	stop   chan struct{}
}

func NewPacketServer(logger Logger) *_PacketServer {
	s := _PacketServer{
		logger: logger,
	}
	if s.logger == nil {
		s.logger = noop.Logger()
	}
	return &s
}

type PacketSender interface {
	SendPacket(Packet) error
}
type PacketSenderFunc func(Packet) error

func (fn PacketSenderFunc) SendPacket(p Packet) error {
	return fn(p)
}

type PacketHandler interface {
	Handle(PacketSender, Packet) error
}

type PacketHandlerFunc func(PacketSender, Packet) error

func (h PacketHandlerFunc) Handle(ps PacketSender, p Packet) error {
	return h(ps, p)
}

func (s *_PacketServer) Shutdown() {
	close(s.stop)
}

func (s *_PacketServer) Serve(listener net.Listener, handler PacketHandler) (serveError error) {
	s.stop = make(chan struct{})

	defer func() {
		rec := recover()
		if rec != nil {
			serveError = errors.New(fmt.Sprint(rec))
		}
	}()

	for {
		select {
		case <-s.stop:
			s.logger.Info("Server stopped gracefully")
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				s.logger.Error("Failed to accept incoming connection",
					"error", err,
				)
				return err
			}
			s.logger.Info("Connection accepted",
				"addr", conn.RemoteAddr().String(),
			)
			go s.handleIncomingConnection(conn, handler)
		}
	}
}

func (s *_PacketServer) handleIncomingConnection(conn net.Conn, handler PacketHandler) {
	defer conn.Close()
	err := s.handleIncomingConnectionE(conn, handler)
	if err != nil && errors.Is(err, io.EOF) {
		s.logger.Error("failed to receive data",
			"error", err,
		)
	}
}

var ErrClientConnClosed = errors.WithMessage(io.EOF, "client connection was closed")

func (s *_PacketServer) handleIncomingConnectionE(in net.Conn, handler PacketHandler) error {
	packet, err := ReceivePacket(in)
	if errors.Is(err, io.EOF) {
		return ErrClientConnClosed
	}
	if err != nil {
		return err
	}

	sender := PacketSenderFunc(func(p Packet) error {
		return SendPacket(in, p)
	})

	err = handler.Handle(sender, packet)
	if err != nil {
		return errors.WithMessage(err, "failed to handle packet")
	}
	return nil
}
