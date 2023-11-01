package app

import (
	"encoding/binary"

	"github.com/brianvoe/gofakeit/v6"
)

func RandomHeader() []byte {
	return binary.BigEndian.AppendUint32(nil, gofakeit.Uint32())
}
