package multi

import (
	"sync"
)

// Degradable starts degraded
type Degradable[T any] struct {
	value    T
	degraded chan struct{}
	lock     sync.RWMutex
}

func (d *Degradable[T]) WhenDegraded() <-chan struct{} {
	if d.degraded == nil {
		ch := make(chan struct{})
		close(ch)
		return ch
	}
	return d.degraded
}

func (d *Degradable[T]) SetDegraded() {
	locked := d.lock.TryLock()
	if !locked {
		return
	}
	defer d.lock.Unlock()
	close(d.degraded)
}

// Value returns the value and a boolean indicating if it is healthy=true or degraded=false
func (d *Degradable[T]) Value() (T, bool) {
	select {
	case <-d.WhenDegraded():
		return d.value, false
	default:
		return d.value, true
	}
}

func (d *Degradable[T]) Reset(value T) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.value = value
	d.degraded = make(chan struct{})
}
