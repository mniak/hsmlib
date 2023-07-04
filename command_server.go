package hsmlib

import (
	"net"
)

type CommandServer struct {
	packetServer PacketServer
}

type CommandHandler interface {
	Handle(req Command) (Response, error)
}

type CommandHandlerFunc func(Command) (Response, error)

func (h CommandHandlerFunc) Handle(c Command) (Response, error) {
	return h(c)
}

func (s *CommandServer) ListenAndServe(address string, handler CommandHandler) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	return s.Serve(listener, handler)
}

func (s *CommandServer) Serve(listener net.Listener, handler CommandHandler) error {
	packetHandler := makePacketHandler(handler)
	return s.packetServer.Serve(listener, packetHandler)
}

func makePacketHandler(cmdHandler CommandHandler) PacketHandler {
	return PacketHandlerFunc(func(ps PacketSender, p Packet) error {
		cmd, err := ParseRawCommand(p.Payload)
		if err != nil {
			return err
		}

		resp, err := cmdHandler.Handle(cmd.WithHeader(p.Header))
		if err != nil {
			return err
		}

		respPacket := Packet{
			Header:  p.Header,
			Payload: resp.serializeResponse(cmd.Code),
		}
		return ps.SendPacket(respPacket)
	})
}
