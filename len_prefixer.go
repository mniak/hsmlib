package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

type PrefixedUint16Protocol struct{}

func (PrefixedUint16Protocol) Send(w io.Writer, data []byte) error {
	var lenbuf [2]byte
	binary.BigEndian.PutUint16(lenbuf[:], uint16(len(data)))

	var buf bytes.Buffer
	_, err := buf.Write(lenbuf[:])
	if err != nil {
		return err
	}

	_, err = buf.Write(data)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, &buf)
	if err != nil {
		return err
	}
	return nil
}

func (PrefixedUint16Protocol) Receive(r io.Reader) ([]byte, error) {
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
