package hsmlib

func ListenAndServeRaw(addr string, handler PacketHandler) error {
	server := PacketServer{}
	return server.ListenAndServe(addr, handler)
}

func ListenAndServe(addr string, handler CommandHandler) error {
	server := CommandServer{}
	return server.ListenAndServe(addr, handler)
}
