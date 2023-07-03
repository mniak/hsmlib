package multi

type Reactor interface {
	Post(data []byte) ([]byte, error)
	Stop()
	Wait()
}
