package multi

import (
	"io"
	"log"
	"time"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/internal/noop"
	"github.com/pkg/errors"
)

type _SimpleReactor struct {
	idManager        IDManager
	target           hsmlib.PacketStream
	logger           hsmlib.Logger
	connectionCloser io.Closer

	used    bool
	stop    chan struct{}
	stopped chan struct{}
	timeout time.Duration
}

func NewSimpleReactor(rw io.ReadWriteCloser) *_SimpleReactor {
	reactor := _SimpleReactor{
		idManager:        &SequentialIDManager{},
		target:           hsmlib.NewPacketStream(rw),
		connectionCloser: rw,
		stop:             make(chan struct{}),
		stopped:          make(chan struct{}),
		timeout:          10 * time.Second,
	}
	return &reactor
}

func (r *_SimpleReactor) Start() error {
	log.Println("Start Begin")
	if r.logger == nil {
		r.logger = noop.Logger()
	}

	if r.used {
		return errors.New("a reactor can only be started once. is was already started before.")
	}
	r.used = true

	go func() {
		defer close(r.stopped)
		defer r.connectionCloser.Close()
		r.handleLoop()
		log.Println("Handle loop stopped. Stopping.")
	}()
	log.Println("Start End")
	return nil
}

func (r *_SimpleReactor) handleLoop() {
	for {
		select {
		case <-r.stop:
			log.Println("Stop signal received")
			return
		default:
			err := r.handleSinglePacket()
			if errors.Is(err, io.EOF) {
				log.Println("Connection closed. stopping.")
				return
			} else if err != nil {
				r.logger.Error("reactor failed and is stopping",
					"error", err,
				)
				log.Println("Failed to handle packet, so stopping:", err)
				return
			}
		}
	}
}

func (r *_SimpleReactor) handleSinglePacket() error {
	log.Println("Receiving packets")
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

var ErrResponseTimeout = errors.New("timeout while waiting for response")

func (r *_SimpleReactor) Post(data []byte) ([]byte, error) {
	timeoutChan := time.After(r.timeout)
	log.Println("-> POST. stopped chan nil?", r.stopped == nil)
	defer log.Println("Post stop")
	select {
	case <-r.stopped:
		log.Println("POST Stopped")
		return nil, errors.New("trying to post into a stopped reactor")
	default:
		log.Println("POST Default")
		id, ch := r.idManager.NewID()
		packet := hsmlib.Packet{
			Header:  id,
			Payload: data,
		}
		err := r.target.SendPacket(packet)
		if err != nil {
			return nil, err
		}
		for {
			select {
			case response := <-ch:
				return response, nil
			case <-timeoutChan:
				return nil, ErrResponseTimeout
			}
		}
	}
}

func (r *_SimpleReactor) Wait() {
	log.Println("Wait started")
	<-r.stopped
	log.Println("Wait finished")
}

func (r *_SimpleReactor) Stop() {
	log.Println("Stopping")
	close(r.stop)
}
