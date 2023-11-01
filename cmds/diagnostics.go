package cmds

func MakeHealthcheck() Healthcheck {
	return Healthcheck{}
}

type Healthcheck struct{}

func (e Healthcheck) Code() []byte {
	return []byte("JK")
}

func (e Healthcheck) Data() []byte {
	return []byte{}
}

type HealthcheckResponse struct{}
