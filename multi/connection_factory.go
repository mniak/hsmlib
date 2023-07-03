package multi

import (
	"io"
	"net"
)

type ConnectionFactory func(target string) (io.ReadWriteCloser, error)

func SimpleTCPDialer() ConnectionFactory {
	return func(target string) (io.ReadWriteCloser, error) {
		return net.Dial("tcp", target)
	}
}

func TCPDialerWithRetry(dialer ConnectionFactory, retryStrategy RetryStrategy) ConnectionFactory {
	return func(target string) (io.ReadWriteCloser, error) {
		return TryAndRetryE(retryStrategy, func() (io.ReadWriteCloser, error, bool) {
			conn, err := dialer(target)
			return conn, err, err != nil
		})
	}
}
