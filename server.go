package hsmlib

import (
	"net"

	"github.com/mniak/hsmlib/adapters/stdlib"
)

type CommandServer struct {
	Logger       Logger
	packetServer PacketServer
}

func (s *CommandServer) Serve(listener net.Listener, handler CommandHandler) error {
	packetHandler := makePacketHandler(handler)
	s.packetServer = NewPacketServer(s.Logger)
	return s.packetServer.Serve(listener, packetHandler)
}

func makePacketHandler(cmdHandler CommandHandler) PacketHandler {
	return PacketHandlerFunc(func(ps PacketSender, p Packet) error {
		cmd, err := ParseCommand(p.Payload)
		if err != nil {
			return err
		}

		resp, err := cmdHandler.Handle(cmd.WithHeader(p.Header))
		if err != nil {
			return err
		}

		respPacket := Packet{
			Header:  p.Header,
			Payload: AddCodeToResponse(resp, cmd.Code()).Bytes(),
		}
		return ps.SendPacket(respPacket)
	})
}

var DefaultLogger Logger = stdlib.NewLogger("[hsmlib] ")

func ListenAndServePackets(addr string, handler PacketHandler) error {
	server := NewPacketServer(DefaultLogger)
	return ListenAndServeI[PacketHandler](server, addr, handler)
}

func ListenAndServe(addr string, handler CommandHandler) error {
	server := CommandServer{
		Logger: DefaultLogger,
	}

	err := ListenAndServeI[CommandHandler](&server, addr, handler)
	return err
}
