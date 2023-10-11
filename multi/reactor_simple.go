package multi

import (
	"io"
	"log"
	"time"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/internal/noop"
	"github.com/pkg/errors"
)

type SimpleReactor struct {
	IDManager IDManager
	Logger    hsmlib.Logger
	Timeout   time.Duration
	Target    hsmlib.PacketStream

	used    bool
	running bool
	stop    chan struct{}
	stopped chan struct{}
}

func (r *SimpleReactor) ensureInit() {
	if r.IDManager == nil {
		r.IDManager = &SequentialIDManager{}
	}
	if r.stop == nil {
		r.stop = make(chan struct{})
	}
	if r.stopped == nil {
		r.stopped = make(chan struct{})
	}
	if r.Timeout == 0 {
		r.Timeout = 10 * time.Second
	}
	if r.Logger == nil {
		r.Logger = noop.Logger()
	}
}

func (r *SimpleReactor) Start() error {
	r.ensureInit()
	if r.Target == nil {
		return errors.New("reactor without target")
	}
	log.Println("Reactor starting")

	if r.used {
		return errors.New("a reactor can only be started once. this has already started before.")
	}
	r.used = true
	r.running = true

	go func() {
		defer close(r.stopped)
		r.handleLoop()
		log.Println("Handle loop stopped. Stopping.")
	}()

	log.Println("Reactor started")
	return nil
}

func (r *SimpleReactor) handleLoop() {
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
				r.Logger.Error("reactor failed and is stopping",
					"error", err,
				)
				log.Println("Failed to handle packet, so stopping:", err)
				return
			}
		}
	}
}

func (r *SimpleReactor) handleSinglePacket() error {
	log.Println("Receiving packets")
	packet, err := r.Target.ReceivePacket()
	if err != nil {
		return errors.WithMessage(err, "could not receive packet")
	}
	channel, found := r.IDManager.FindChannel(packet.Header)
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

func (r *SimpleReactor) Post(data []byte) ([]byte, error) {
	if !r.running {
		return nil, errors.New("cant post in reactor that is not running")
	}

	timeoutChan := time.After(r.Timeout)
	log.Println("-> POST. stopped chan nil?", r.stopped == nil)
	defer log.Println("Post stop")
	select {
	case <-r.stopped:
		log.Println("POST Stopped")
		return nil, errors.New("trying to post into a stopped reactor")
	default:
		log.Println("POST Default")
		id, ch := r.IDManager.NewID()
		packet := hsmlib.Packet{
			Header:  id,
			Payload: data,
		}
		err := r.Target.SendPacket(packet)
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

func (r *SimpleReactor) Wait() {
	log.Println("Wait started")
	<-r.stopped
	log.Println("Wait finished")
}

func (r *SimpleReactor) Stop() {
	log.Println("Stopping")
	close(r.stop)
}
