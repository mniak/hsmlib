package multi

import (
	"encoding/binary"
	"sync"
	"sync/atomic"

	"github.com/mniak/hsmlib"
)

type IDManager interface {
	IDLength() int
	NewID() ([]byte, <-chan []byte)
	FindChannel(id []byte) (chan<- []byte, bool)
	CloseAllChannels()
}

type SequentialIDManager struct {
	ids     sync.Map
	counter atomic.Uint32
}

func (m *SequentialIDManager) IDLength() int {
	return hsmlib.HeaderLength
}

func (m *SequentialIDManager) NewID() ([]byte, <-chan []byte) {
	newID := m.counter.Add(1)
	callbackChanI, _ := m.ids.LoadOrStore(newID, make(chan []byte))
	callbackChan := callbackChanI.(chan []byte)
	newIDBytes := binary.BigEndian.AppendUint32(nil, newID)
	return newIDBytes, callbackChan
}

func (m *SequentialIDManager) FindChannel(id []byte) (chan<- []byte, bool) {
	idnum := binary.BigEndian.Uint32(id)
	callbackChanI, found := m.ids.Load(idnum)
	if !found {
		return nil, false
	}

	callbackChan := callbackChanI.(chan []byte)
	return callbackChan, true
}

func (m *SequentialIDManager) CloseAllChannels() {
	oldMap := &m.ids
	m.ids = sync.Map{}
	oldMap.Range(func(key, value any) bool {
		channel := value.(chan []byte)
		close(channel)
		return true
	})
}
