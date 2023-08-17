package multi

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/mniak/hsmlib"
	"github.com/mniak/hsmlib/internal/noop"
	"github.com/pkg/errors"
)

type _ResilientReactor struct {
	target string
	logger hsmlib.Logger

	inner         Degradable[Reactor]
	runLock       sync.Mutex
	retryStrategy RetryStrategy
	stop          chan struct{}
	stopped       chan struct{}
	ractorFactory ReactorFactory
}

func NewResilientReactor(target string, logger hsmlib.Logger) *_ResilientReactor {
	r := _ResilientReactor{
		target:        target,
		retryStrategy: SimpleDelayRetryStrategy(1 * time.Second),
		logger:        logger,
	}
	if r.logger == nil {
		r.logger = noop.Logger()
	}
	return &r
}

func (r *_ResilientReactor) Start() error {
	canRun := r.runLock.TryLock()
	if !canRun {
		return errors.New("reactor is already running")
	}
	if err := r.connectAndSaveReactor(false); err != nil {
		return errors.WithMessage(err, "first connection did not work")
	}
	go func() {
		defer r.runLock.Unlock()
		r.reconnectLoop()
	}()
	return nil
}

func (r *_ResilientReactor) connectAndSaveReactor(withRetry bool) error {
	if inner, _ := r.inner.Value(); inner != nil {
		inner.Stop()
	}

	reactor, err := r.ractorFactory(r.target, withRetry)
	if err != nil {
		return err
	}
	r.logger.Info("connect() worked. resetting degradable.")

	go func() {
		reactor.Wait()
		r.inner.SetDegraded()
	}()
	r.inner.Reset(reactor)
	return nil
}

func (r *_ResilientReactor) reconnectLoop() {
	for {
		select {
		case <-r.stop:
			return
		case <-r.inner.WhenDegraded():
			r.logger.Error("detected degraded reactor. trying to reconnect.")
			err := r.connectAndSaveReactor(true)
			if err != nil {
				r.logger.Error("failed to reconnect",
					"error", err,
				)
			}
		}
	}
}

var ReactorDegraded = errors.New("reactor is degraded")

func (r *_ResilientReactor) Post(data []byte) ([]byte, error) {
	reactor, healty := r.inner.Value()
	if !healty {
		r.logger.Error("not trying to post because reactor is degraded")
		return nil, ReactorDegraded
	}
	resp, err := reactor.Post(data)
	switch {
	case err == nil:
		return resp, err
	case errors.Is(err, io.EOF), errors.Is(err, net.ErrClosed):
		r.logger.Error("post failed. marking reactor as degraded",
			"error", err,
		)
		r.inner.SetDegraded()
		return resp, ReactorDegraded
	default:
		r.logger.Error("post failed",
			"error", err,
		)
		return resp, err
	}
}

func (r *_ResilientReactor) Wait() {
	<-r.stopped
}

func (r *_ResilientReactor) Stop() {
	close(r.stop)
}
