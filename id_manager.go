package main

import (
	"encoding/binary"
	"sync"
	"sync/atomic"
)

type IDManager interface {
	IDLength() int
	NewID() ([]byte, <-chan []byte)
	FindChannel(id []byte) (chan<- []byte, bool)
}

type SequentialUint32IDManager struct {
	ids     sync.Map
	counter atomic.Uint32
}

func (m *SequentialUint32IDManager) IDLength() int {
	return 4
}

func (m *SequentialUint32IDManager) NewID() ([]byte, <-chan []byte) {
	newID := m.counter.Add(1)
	callbackChanI, _ := m.ids.LoadOrStore(newID, make(chan []byte))
	callbackChan := callbackChanI.(chan []byte)
	newIDBytes := binary.BigEndian.AppendUint32(nil, newID)
	return newIDBytes, callbackChan
}

func (m *SequentialUint32IDManager) FindChannel(id []byte) (chan<- []byte, bool) {
	callbackChanI, found := m.ids.Load(id)
	if !found {
		return nil, false
	}

	callbackChan := callbackChanI.(chan []byte)
	return callbackChan, true
}
