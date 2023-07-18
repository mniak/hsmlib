package hsmlib

import (
	"encoding/binary"
	"io"
)

func ReceiveFrame(r io.Reader) ([]byte, error) {
	var lenbuf [2]byte
	_, err := io.ReadFull(r, lenbuf[:])
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(lenbuf[:])
	databuf := make([]byte, length)
	_, err = io.ReadFull(r, databuf)
	if err != nil {
		return nil, err
	}
	return databuf, nil
}

func SendFrame(w io.Writer, payload []byte) error {
	result := LengthPrefix2B(payload)
	_, err := w.Write(result)
	return err
}
