package hsmlib

import (
	"encoding/binary"
	"io"
)

func ReceiveFrame(r io.Reader) ([]byte, error) {
	var lenbuf [2]byte
	_, err := r.Read(lenbuf[:])
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(lenbuf[:])
	databuf := make([]byte, length)
	_, err = r.Read(databuf)
	if err != nil {
		return nil, err
	}
	return databuf, nil
}

func SendFrame(w io.Writer, payload []byte) error {
	var lenbuf [2]byte
	binary.BigEndian.PutUint16(lenbuf[:], uint16(len(payload)))

	result := append(lenbuf[:], payload...)

	_, err := w.Write(result)
	return err
}
