package nsq

import (
	"github.com/deliveryhero/pipeline"
	"time"

	"github.com/fatih/color"
	"github.com/segmentio/nsq-go"

	"github.com/vortex14/gotyphoon/interfaces"
)

type Consumer struct {
	total int
	Name string
	Topic string
	Concurrent int
	Channel string
	Worker *nsq.Consumer
}

func (c *Consumer) GetTotal() int {
	return c.total
}

func (s *Service) InitConsumer(
	settings *interfaces.Queue,
) *nsq.Consumer {

	name := settings.GetGroupName()
	priority := settings.GetPriority()

	color.Blue(`init NSQ consumer. 
	Group:%s 
	Topic: %s 
	Channel: %s`,  name, settings.Topic, settings.Channel,
	)

	s.initConfig()
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:   settings.Topic,
		Channel: settings.Channel,
		Address: s.Config.NsqdNodes[0].IP,
		ReadTimeout: time.Duration(60) * time.Second,
		MaxInFlight: settings.Concurrent,
	})

	newConsumer := &Consumer{
		Topic:      settings.Topic,
		Channel:    settings.Channel,
		Concurrent: settings.Concurrent,
		Worker:     consumer,
	}

	if s.Consumers == nil {
		s.Consumers = Consumers{
			name: {
				newConsumer,
			},
		}
	} else {
		s.Consumers[name] = append(s.Consumers[name], newConsumer)
	}
	color.Yellow("PRIORITY ------ >>> %d", priority)
	if priority > 0 {

		if s.priorityConsumers[name] != nil {
			s.priorityConsumers[name][priority] = newConsumer
		} else {
			s.priorityConsumers[name] = map[int]*Consumer{
				priority: newConsumer,
			}
		}

	}





	return consumer
}

func (s *Service) StopConsumers()  {
	for name := range s.Consumers {
		for _, consumer := range s.Consumers[name] {
			color.Yellow("stop nsq consumer: %s", name)
			consumer.Worker.Stop()
		}
	}
}


func (s *Service) read(consumer *Consumer) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for msg := range consumer.Worker.Messages() {
			out <- &interfaces.YieldNsqMessage{
				Msg:     &msg,
				Yield: &interfaces.Yield{
					Name: consumer.Name,
					Topic:   consumer.Topic,
					Channel: consumer.Channel,
				},
			}
		}
	}()

	return out
}



func (s *Service) Read() <-chan *interfaces.YieldNsqMessage {
	out := make(chan *interfaces.YieldNsqMessage)
	go func() {
		defer close(out)
		for source := range pipeline.Merge(s.mergeConsumerChannels()...) {
			out <- source.(*interfaces.YieldNsqMessage)
		}
	}()
	return out
}

