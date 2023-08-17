package multi

import (
	"time"

	"github.com/pkg/errors"
)

type ReactorFactory func(target string, withRetry bool) (Reactor, error)

func TCPReactorFactory() ReactorFactory {
	getDialer := func(withRetry bool) ConnectionFactory {
		dialer := SimpleTCPDialer()
		if withRetry {
			dialer = TCPDialerWithRetry(dialer, SimpleDelayRetryStrategy(1*time.Second))
		}
		return dialer
	}
	return func(target string, withRetry bool) (Reactor, error) {
		conn, err := getDialer(withRetry)(target)
		if err != nil {
			return nil, errors.WithMessage(err, "could not connect")
		}
		reactor := NewSimpleReactor(conn)
		if err := reactor.Start(); err != nil {
			return reactor, errors.WithMessage(err, "could not start reactor")
		}
		return reactor, err
	}
}
