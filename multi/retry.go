package multi

import (
	"time"

	"github.com/pkg/errors"
)

var ErrMustRetry = errors.New("must retry")

func TryAndRetryE[T any](strategy RetryStrategy, do func() (result T, err error, mustRetry bool)) (result T, err error) {
	strategy.Try(func() bool {
		var mustRetry bool
		result, err, mustRetry = do()
		return mustRetry
	})
	return result, err
}

func TryAndRetry[T any](strategy RetryStrategy, do func() (result T, mustRetry bool)) (result T) {
	strategy.Try(func() bool {
		var mustRetry bool
		result, mustRetry = do()
		return mustRetry
	})
	return result
}

type _Try interface{}

type RetryStrategy interface {
	Try(do func() bool)
}

type inlineRetryStrategy func(do func() bool)

func (strat inlineRetryStrategy) Try(do func() bool) {
	strat(do)
}

func SimpleDelayRetryStrategy(delay time.Duration) RetryStrategy {
	return inlineRetryStrategy(func(do func() bool) {
		for {
			mustRetry := do()
			if mustRetry {
				time.Sleep(delay)
			} else {
				break
			}
		}
	})
}
