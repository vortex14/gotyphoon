package nsq

import (
	"github.com/fatih/color"
	"github.com/segmentio/nsq-go"

	"github.com/vortex14/gotyphoon/interfaces"
)

type Producer struct {
	total  int
	Worker *nsq.Producer
}

func (s *Service) PriorityPub(priority int, name string, message string) {
	err := s.priorityProducers[name][priority].Worker.Publish([]byte(message))
	if err != nil {
		color.Red("%s", err.Error())
	}
}

func (p *Producer) GetTotal() int {
	return p.total
}

func (s *Service) Pub(
	name string,
	topic string,
	message string,
) error {
	return s.Producers[name].Worker.PublishTo(topic, []byte(message))
}

func (s *Service) InitProducer(settings *interfaces.Queue) {
	name := settings.GetGroupName()
	priority := settings.GetPriority()
	color.Yellow(`init NSQ producer. 
	Group:%s 
	Topic: %s 
	Channel: %s`, name, settings.Topic, settings.Channel,
	)
	s.initConfig()
	producer, _ := nsq.StartProducer(nsq.ProducerConfig{
		Topic:   settings.Topic,
		Address: s.Config.NsqdNodes[0].IP,
	})

	newProducer := &Producer{Worker: producer}

	if s.Producers == nil {
		s.Producers = Producers{
			name: newProducer,
		}
	} else {
		s.Producers[name] = newProducer
	}

	if _, ok := s.Producers[name]; !ok {
		panic("producer cannot be installed on the map ")
	}

	if priority > 0 {

		if s.priorityProducers[name] != nil {
			s.priorityProducers[name][priority] = newProducer
		} else {
			s.priorityProducers[name] = map[int]*Producer{
				priority: newProducer,
			}
		}

	}

}

func (s *Service) StopProducers() {
	color.Yellow("Stop all producers ...")
	for producer := range s.Producers {
		s.Producers[producer].Worker.Stop()
	}

}
