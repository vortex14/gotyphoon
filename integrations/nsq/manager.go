package nsq

import (
	"github.com/vortex14/gotyphoon/interfaces"
	"strings"

	"github.com/fatih/color"
	"github.com/segmentio/nsq-go"
)

func (s *Service) StopNSQ()  {
	s.StopProducers()
	s.StopConsumers()
}

func (s *Service) Ping() bool  {
	status := false
	s.initConfig()
	NsqlookupdIP := strings.ReplaceAll(s.Config.NsqlookupdIP, "http://","")
	client := nsq.Client{Address: NsqlookupdIP}
	err := client.Ping()
	if err != nil {
		color.Red("%s", err)
	} else {
		status = true
	}

	return status
}


func (s *Service) GetHost() string {
	return s.Config.NsqlookupdIP
}

func (s *Service) SetOptions(options interfaces.MessageBrokerOptions)  {
	s.Options = options
}


func (s *Service) GetPort() int {

	//TODO:
	return PORT
}
