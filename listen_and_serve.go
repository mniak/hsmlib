package hsmlib

import "net"

type (
	ServeFunc[T any] func(listener net.Listener, handler T) error
	ServeI[T any]    interface {
		Serve(listener net.Listener, t T) error
	}
)

func ListenAndServeFn[T any](serve ServeFunc[T], address string, handler T) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	return serve(listener, handler)
}

func ListenAndServeI[T any](server ServeI[T], address string, t T) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	return server.Serve(listener, t)
}
