package hsmlib

import (
	"net"

	"github.com/mniak/hsmlib/internal/noop"
	"github.com/pkg/errors"
)

type PacketServer struct {
	Logger Logger
	stop   chan struct{}
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

func (s *PacketServer) ListenAndServe(address string, handler PacketHandler) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	return s.Serve(listener, handler)
}

func (s *PacketServer) Shutdown() {
	close(s.stop)
}

func (s *PacketServer) Serve(listener net.Listener, handler PacketHandler) error {
	if s.Logger == nil {
		s.Logger = noop.Logger()
	}

	s.stop = make(chan struct{})

	for {
		select {
		case <-s.stop:
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				s.Logger.Error("Failed to accept incoming connection",
					"error", err,
				)
				return err
			}
			s.Logger.Info("Connection accepted",
				"addr", conn.RemoteAddr().String(),
			)
			go s.handleIncomingConnection(conn, handler)
		}
	}
}

func (s *PacketServer) handleIncomingConnection(conn net.Conn, handler PacketHandler) {
	err := s.handleIncomingConnectionE(conn, handler)
	if err != nil {
		s.Logger.Error("failed to receive data",
			"error", err,
		)
	}
}

func (s *PacketServer) handleIncomingConnectionE(conn net.Conn, handler PacketHandler) error {
	defer conn.Close()

	packet, err := ReceivePacket(conn)
	if err != nil {
		return err
	}

	sender := PacketSenderFunc(func(p Packet) error {
		return SendPacket(conn, p)
	})

	err = handler.Handle(sender, packet)
	if err != nil {
		return errors.WithMessage(err, "failed to handle packet")
	}
	return nil
}
