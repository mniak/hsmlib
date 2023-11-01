package cmds

import (
	"errors"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/errcode"
)

func MakeEcho(message string) Echo {
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

func EchoFromCommand(cmd hsmlib.Command) (Echo, error) {
	msg, remaining, err := hsmlib.ParseWithLengthPrefix4H(cmd.Data())
	if err != nil || len(remaining) > 0 {
		return Echo{}, errors.New("echo message length on header does not correspond to actual length")
	}

	return MakeEcho(string(msg)), nil
}

type EchoResponse struct {
	ErrCode errcode.ErrorCode
	Message string
}

func (r EchoResponse) ErrorCode() errcode.ErrorCode {
	return r.ErrCode
}

func (r EchoResponse) Data() []byte {
	return []byte(r.Message)
}

func MakeEchoResponse(ec errcode.ErrorCode, msg string) EchoResponse {
	return EchoResponse{
		ErrCode: ec,
		Message: msg,
	}
}
