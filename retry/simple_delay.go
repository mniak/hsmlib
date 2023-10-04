package retry

import "time"

func SimpleDelayStrategy(delay time.Duration) Strategy {
	return inlineStrategy(func(do func() (success bool)) {
		for {
			success := do()
			if !success {
				time.Sleep(delay)
			} else {
				break
			}
		}
	})
}
