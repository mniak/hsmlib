package multi

import (
	"context"
	"io"
	"net"
	"sync"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/internal/noop"
	"github.com/pkg/errors"
)

type Reactor interface {
	Post(ctx context.Context, data []byte) ([]byte, error)
}

type _Reactor struct {
	idManager IDManager
	target    hsmlib.PacketStream
	logger    hsmlib.Logger
	closer    io.Closer

	run     sync.Mutex
	stop    chan struct{}
	stopped chan struct{}
}

func NewReactorFromReadWriter(rw io.ReadWriteCloser) *_Reactor {
	reactor := _Reactor{
		idManager: &SequentialIDManager{},
		target:    hsmlib.NewPacketStream(rw),
		closer:    rw,
	}
	return &reactor
}

func NewReactor(target string) (*_Reactor, error) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		return nil, err
	}
	return NewReactorFromReadWriter(conn), nil
}

func (r *_Reactor) Start() error {
	if r.logger == nil {
		r.logger = noop.Logger()
	}

	canRun := r.run.TryLock()
	if !canRun {
		return errors.New("reactor is already running")
	}
	defer r.run.Unlock()

	r.stop = make(chan struct{})
	r.stopped = make(chan struct{})
	go func() {
		defer close(r.stopped)
		defer r.closer.Close()
		for {
			select {
			case <-r.stop:
				return
			default:
				err := r.handleOnPacket()
				if err != nil {
					r.logger.Error("reactor failed and is stopping",
						"error", err,
					)
					return
				}
			}
		}
	}()
	return nil
}

func (r *_Reactor) handleOnPacket() error {
	packet, err := r.target.ReceivePacket()
	if err != nil {
		return errors.WithMessage(err, "could not receive packet")
	}
	channel, found := r.idManager.FindChannel(packet.Header)
	if !found {
		return errors.WithMessagef(err, "callback channel '%2X' not found", packet.Header)
	}

	go func() {
		channel <- packet.Payload
		close(channel)
	}()
	return nil
}

func (r *_Reactor) Post(ctx context.Context, data []byte) ([]byte, error) {
	id, ch := r.idManager.NewID()
	packet := hsmlib.Packet{
		Header:  id,
		Payload: data,
	}
	err := r.target.SendPacket(packet)
	if err != nil {
		return nil, err
	}

	select {
	case response := <-ch:
		return response, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *_Reactor) Wait() {
	<-r.stopped
}

func (r *_Reactor) Stop() {
	close(r.stop)
}
