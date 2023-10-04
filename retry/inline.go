package retry

type inlineStrategy func(do func() (success bool))

func (s inlineStrategy) Try(do func() (success bool)) {
	s(do)
}
