package hsmlib

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
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

func ParseWithLengthPrefix4H(data []byte) (read []byte, remaining []byte, err error) {
	if len(data) < 4 {
		return nil, nil, errors.New("could not decode length 4H: too short")
	}
	lengthHexBytes := data[:4]
	data = data[4:]

	lengthBytes, err := hex.DecodeString(string(lengthHexBytes))
	if err != nil {
		return nil, nil, errors.New("could not decode length 4H: invalid characters")
	}

	length := binary.BigEndian.Uint16(lengthBytes)

	if len(data) < int(length) {
		return nil, nil, errors.New("could not decode length 4H: length bigger than remaining data")
	}

	return data[:length], data[length:], nil
}
