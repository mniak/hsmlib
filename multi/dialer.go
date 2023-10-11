package multi

import (
	"context"
	"crypto/tls"
	"io"
	"net"

	"github.com/mniak/hsmlib/retry"
)

type Dialer interface {
	DialContext(ctx context.Context) (io.ReadWriteCloser, error)
}

type TCPDialer struct {
	Address string
}

func (self TCPDialer) DialContext(ctx context.Context) (io.ReadWriteCloser, error) {
	d := net.Dialer{}
	return d.DialContext(ctx, "tcp", self.Address)
}

type TLSDialer struct {
	Address   string
	TLSConfig *tls.Config
}

func (self TLSDialer) DialContext(ctx context.Context) (io.ReadWriteCloser, error) {
	d := tls.Dialer{
		Config: self.TLSConfig,
	}
	return d.DialContext(ctx, "tcp", self.Address)
}

type RetryDialer struct {
	InnerDialer   Dialer
	RetryStrategy retry.Strategy
}

func (self RetryDialer) Dial(ctx context.Context) (io.ReadWriteCloser, error) {
	return retry.TryE(self.RetryStrategy, func() (io.ReadWriteCloser, error, bool) {
		conn, err := self.InnerDialer.DialContext(ctx)
		if err == context.DeadlineExceeded {
			return nil, err, false
		}
		return conn, err, err != nil
	})
}
