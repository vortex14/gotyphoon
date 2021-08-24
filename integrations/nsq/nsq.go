package nsq

import (
	"github.com/sirupsen/logrus"
	"strconv"

	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	PORT = 4161
)

type Producers map[string] *Producer
type Consumers map[string] [] *Consumer


type Service struct {
	Producers Producers
	Consumers Consumers
	Project interfaces.Project
	Config *interfaces.ConfigProject
	priorityProducers map[string]map[int] *Producer
	priorityConsumers map[string]map[int] *Consumer
	Options interfaces.MessageBrokerOptions
}


func (s *Service) RunNSQ()  {

	logrus.Debug("Running connections to NSQ ...")

	s.initConfig()
	for qSettings := range s.getSettings() {
		group := qSettings.GetGroupName()
		if group == interfaces.PRIORITY ||
			group == interfaces.PROCESSOR2PRIORITY ||
			group == interfaces.TRANSPORTER2PRIORITY {
			//			color.Yellow(`
			//	Init queue group: %s
			//	Topic:			  %s
			//	Channel: 		  %s
			//	W/R:			  %t/%t
			//`, group, qSettings.Topic, qSettings.Channel, qSettings.Writable, qSettings.Readable)
			for _, i := range []int{1,2,3} {
				prioritySetting := &interfaces.Queue{
					Concurrent: qSettings.Concurrent,
					MsgTimeout: qSettings.MsgTimeout,
					Channel:    qSettings.Channel,
					Topic:      s.Config.GetTopic(qSettings.GetComponentName(), group, strconv.Itoa(i)),
					Share:      qSettings.Share,
					Writable:   qSettings.Writable,
					Readable:   qSettings.Readable,
				}
				prioritySetting.SetPriority(i)
				prioritySetting.SetGroupName(group)
				s.InitQueue(prioritySetting)
			}

			continue
		}
		qSettings.Topic = s.Config.GetTopic(qSettings.GetComponentName(), group, "")
		s.InitQueue(qSettings)
	}
}

