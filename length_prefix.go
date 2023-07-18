package hsmlib

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func LengthPrefix2B(data []byte) []byte {
	var lenbuf [2]byte
	binary.BigEndian.PutUint16(lenbuf[:], uint16(len(data)))
	result := append(lenbuf[:], data...)
	return result
}

func LeftPad(s, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

func LengthPrefix4H(data []byte) []byte {
	prefix := strings.ToUpper(fmt.Sprintf("%04X", len(data)))
	result := append([]byte(prefix), data...)
	return result
}
