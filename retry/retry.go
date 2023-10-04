package retry

type Strategy interface {
	Try(do func() (success bool))
}

func TryE[T any](strategy Strategy, do func() (result T, err error, success bool)) (T, error) {
	type ResultWithError struct {
		Result T
		Error  error
	}
	result := Try[ResultWithError](strategy, func() (result ResultWithError, success bool) {
		r, e, s := do()
		return ResultWithError{
			Result: r,
			Error:  e,
		}, s
	})
	return result.Result, result.Error
}

func Try[T any](strategy Strategy, do func() (result T, success bool)) T {
	var result T
	strategy.Try(func() bool {
		var success bool
		result, success = do()
		return success
	})
	return result
}
