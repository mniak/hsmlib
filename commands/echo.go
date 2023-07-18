package commands

import "github.com/mniak/hsmlib"

func Echo(message string) hsmlib.Command {
	return hsmlib.Command{
		Code: "B2",
		Data: hsmlib.LengthPrefix4H([]byte(message)),
	}
}
