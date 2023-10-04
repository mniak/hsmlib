package commands

import "github.com/mniak/hsmlib"

func MakeEcho(message string) hsmlib.Command {
	return Echo{
		Message: message,
	}
}

type Echo struct {
	Message string
}

func (e Echo) Code() []byte {
	return []byte("B2")
}

func (e Echo) Data() []byte {
	return hsmlib.LengthPrefix4H([]byte(e.Message))
}

type EchoResponse struct{}
