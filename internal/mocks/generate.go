package mocks

import (
	_ "go.uber.org/mock/mockgen/model"
)

//go:generate mockgen -package=mocks -destination=mocks_io.go io Closer
//go:generate mockgen -package=mocks -destination=mocks_hsmlib.go github.com/mniak/hsmlib PacketStream,Logger
//go:generate mockgen -package=mocks -destination=mocks_multi.go github.com/mniak/hsmlib/multi IDManager
