package main

import "io"

type Protocol[T any] interface {
	Receive(r io.Reader) (T, error)
	Send(w io.Writer, data T) error
}

func NewHSMProtocol() SimpleIDWrapper {
	lenPrefix := PrefixedUint16Protocol{}
	idMan := SequentialUint32IDManager{}
	idWrap := SimpleIDWrapper{
		IDLength: idMan.IDLength(),
		Inner:    lenPrefix,
	}
	return idWrap
}
