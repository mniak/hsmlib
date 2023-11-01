package errcode

import "fmt"

type ErrorCode byte

func (code ErrorCode) Code() string {
	return fmt.Sprintf("%02X", byte(code))
}

func (code ErrorCode) Short() string {
	desc, found := shortDescriptions[code]
	if !found {
		return "Unknown"
	}
	return desc
}

func (code ErrorCode) Long() string {
	desc, found := longDescriptions[code]
	if !found {
		return "Unknown"
	}
	return desc
}

func (code ErrorCode) String() string {
	return fmt.Sprintf("[%s] %s", code.Code(), code.Short())
}
