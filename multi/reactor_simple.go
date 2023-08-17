package multi

import (
	"context"
	"io"
	"sync"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/internal/noop"
	"github.com/pkg/errors"
)

type _SimpleReactor struct {
	idManager        IDManager
	target           hsmlib.PacketStream
	logger           hsmlib.Logger
	connectionCloser io.Closer

	runLock sync.Mutex
	stop    chan struct{}
	stopped chan struct{}
}

func NewSimpleReactor(rw io.ReadWriteCloser) *_SimpleReactor {
	reactor := _SimpleReactor{
		idManager:        &SequentialIDManager{},
		target:           hsmlib.NewPacketStream(rw),
		connectionCloser: rw,
		stop:             make(chan struct{}),
		stopped:          make(chan struct{}),
	}
	return &reactor
}

func (r *_SimpleReactor) Start() error {
	if r.logger == nil {
		r.logger = noop.Logger()
	}

	canRun := r.runLock.TryLock()
	if !canRun {
		return errors.New("reactor is already running")
	}
	defer r.runLock.Unlock()

	r.stop = make(chan struct{})
	r.stopped = make(chan struct{})
	go func() {
		defer close(r.stopped)
		defer r.connectionCloser.Close()
		r.handleLoop()
	}()
	return nil
}

func (r *_SimpleReactor) handleLoop() {
	for {
		select {
		case <-r.stop:
			return
		default:
			err := r.handleSinglePacket()
			if err != nil {
				r.logger.Error("reactor failed and is stopping",
					"error", err,
				)
				return
			}
		}
	}
}

func (r *_SimpleReactor) handleSinglePacket() error {
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

func (r *_SimpleReactor) Post(ctx context.Context, data []byte) ([]byte, error) {
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
	case <-r.stopped:
		return nil, errors.New("trying to post into a stopped reactor")
	case response := <-ch:
		return response, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *_SimpleReactor) Wait() {
	<-r.stopped
}

func (r *_SimpleReactor) Stop() {
	close(r.stop)
}
