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

// var ErrClientConnClosed = errors.WithMessage(io.EOF, "client connection was closed")

func (s *_PacketServer) handleIncomingConnection(conn io.ReadWriteCloser, handler PacketHandler) {
	defer func() {
		conn.Close()
		s.logger.Info("Incoming client connection closed")
	}()
	err := s.handleIncomingConnectionE(conn, handler)
	if err != nil {
		s.logger.Error("failed to receive data",
			"error", err,
		)
	}
}

func (s *_PacketServer) handleIncomingConnectionE(conn io.ReadWriteCloser, handler PacketHandler) error {
	packetStream := NewPacketStream(conn)
	packet, err := packetStream.ReceivePacket()
	if err != nil {
		return err
	}

	err = handler.Handle(packetStream, packet)
	if err != nil {
		return errors.WithMessage(err, "failed to handle packet")
	}
	return nil
}
