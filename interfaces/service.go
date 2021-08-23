package interfaces



type Service interface {
	GetHost() string
	GetPort() int
	Ping() bool
}

type AdapterService interface {
	Ping() bool
	Init()
}

func (s *ServiceRedis) GetHost() string  {
	return s.Details.Host
}

func (s *ServiceRedis) GetPort() int {
	return s.Details.Port
}




func (s *ServiceMongo) GetHost() string  {
	return s.Details.Host
}

func (s *ServiceMongo) GetPort() int {
	return s.Details.Port
}
