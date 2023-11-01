package multi

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/internal/noop"
	"github.com/mniak/hsmlib/retry"
	"github.com/pkg/errors"
)

type ResilientReactor struct {
	Logger            hsmlib.Logger
	RetryStrategy     retry.Strategy
	ConnectionTimeout time.Duration

	degradableReactor Degradable[Reactor]
	runLock           sync.Mutex
	stop              chan struct{}
	stopped           chan struct{}
}

func (r *ResilientReactor) ensureInit() {
	r.stop = make(chan struct{})
	r.stop = make(chan struct{})

	if r.Logger == nil {
		r.Logger = noop.Logger()
	}
	if r.RetryStrategy == nil {
		r.RetryStrategy = retry.SimpleDelayStrategy(1 * time.Second)
	}
}

func (r *ResilientReactor) Start(dialer Dialer) error {
	r.ensureInit()

	canRun := r.runLock.TryLock()
	if !canRun {
		return errors.New("reactor is already running")
	}
	if err := r.connectAndSaveReactor(dialer, false); err != nil {
		return errors.WithMessage(err, "first connection did not work")
	}
	go func() {
		defer r.runLock.Unlock()
		r.reconnectLoop(dialer)
	}()
	return nil
}

func (r *ResilientReactor) connectAndSaveReactor(dialer Dialer, withRetry bool) error {
	if inner, _ := r.degradableReactor.Value(); inner != nil {
		inner.Stop()
	}

	ctx, _ := context.WithTimeout(context.Background(), DefaultTimeout)
	conn, err := dialer.DialContext(ctx)
	if err != nil {
		return err
	}
	reactor := SimpleReactor{
		Target: hsmlib.NewPacketStream(conn).WithLogger(r.Logger),
		Logger: r.Logger,
	}
	err = reactor.Start()
	if err != nil {
		return err
	}
	r.Logger.Info("connect() worked. resetting degradable.")

	go func() {
		reactor.Wait()
		r.degradableReactor.SetDegraded()
	}()
	r.degradableReactor.Reset(&reactor)
	return nil
}

func (r *ResilientReactor) reconnectLoop(dialer Dialer) {
	for {
		select {
		case <-r.stop:
			return
		case <-r.degradableReactor.WhenDegraded():
			r.Logger.Error("detected degraded reactor. trying to reconnect.")
			err := r.connectAndSaveReactor(dialer, true)
			if err != nil {
				r.Logger.Error("failed to reconnect",
					"error", err,
				)
			}
		}
	}
}

var ReactorDegraded = errors.New("reactor is degraded")

func (self *ResilientReactor) Post(data []byte) ([]byte, error) {
	reactor, healty := self.degradableReactor.Value()
	if !healty {
		self.Logger.Error("not trying to post because reactor is degraded")
		return nil, ReactorDegraded
	}
	resp, err := reactor.Post(data)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, io.EOF), errors.Is(err, net.ErrClosed):
		self.Logger.Error("post failed. marking reactor as degraded",
			"error", err,
		)
		self.degradableReactor.SetDegraded()
		return resp, ReactorDegraded
	default:
		self.Logger.Error("post failed",
			"error", err,
		)
		return resp, err
	}
}

func (self *ResilientReactor) Wait() {
	<-self.stopped
}

func (self *ResilientReactor) Stop() {
	close(self.stop)
}
