package nsq

import (
	"github.com/fatih/color"
	"github.com/segmentio/nsq-go"
	"github.com/vortex14/gotyphoon/config"
	"strings"
)

type Service struct {
	Config *config.Config
	Producers map[string] *nsq.Producer
	Consumers map[string] []*nsq.Consumer
}

func (s *Service) TestConnect() bool  {
	status := false
	projectConfig := s.Config
	NsqlookupdIP := strings.ReplaceAll(projectConfig.NsqlookupdIP, "http://","")
	client := nsq.Client{Address: NsqlookupdIP}
	err := client.Ping()
	if err != nil {
		color.Red("%s", err)
	} else {
		status = true
	}

	return status
}

func (s *Service) Pub()  {

}

func (s *Service) Read()  {

}

