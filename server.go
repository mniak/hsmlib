package hsmlib

import "net"

type CommandServer struct {
	packetServer PacketServer
}

type CommandHandler interface {
	Handle(cmd CommandWithHeader) (Response, error)
}

type CommandHandlerFunc func(cmd CommandWithHeader) (Response, error)

func (h CommandHandlerFunc) Handle(cmd CommandWithHeader) (Response, error) {
	return h(cmd)
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
			Payload: resp.WithCode(cmd.Code).Bytes(),
		}
		return ps.SendPacket(respPacket)
	})
}

func ListenAndServeRaw(addr string, handler PacketHandler) error {
	server := PacketServer{}
	return server.ListenAndServe(addr, handler)
}

func ListenAndServe(addr string, handler CommandHandler) error {
	server := CommandServer{}
	return server.ListenAndServe(addr, handler)
}
