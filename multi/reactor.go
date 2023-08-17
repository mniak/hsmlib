package multi

import (
	"context"
)

type Reactor interface {
	Post(ctx context.Context, data []byte) ([]byte, error)
	Stop()
	Wait()
}
