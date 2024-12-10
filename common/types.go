package common

type Handler interface {
	HandleService()
}

type Service interface {
	GetServiceName() string
}
