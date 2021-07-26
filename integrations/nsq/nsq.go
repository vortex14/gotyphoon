package nsq

import (
	"context"
	"github.com/deliveryhero/pipeline"
	"github.com/fatih/color"
	"github.com/segmentio/nsq-go"
	"github.com/vortex14/gotyphoon/interfaces"
	"strings"
	"time"
)

type Consumer struct {
	total int
	Name string
	Topic string
	Concurrent int
	Channel string
	Worker *nsq.Consumer
}


type Producer struct {
	total int
	Worker *nsq.Producer
}

type Producers map[string] *Producer
type Consumers map[string] [] *Consumer


type Service struct {
	Config *interfaces.ConfigProject
	Project interfaces.Project
	Producers Producers
	Consumers Consumers
}

func (s *Service) initConfig()  {
	if s.Config == nil {
		s.Config = s.Project.LoadConfig()
	}
}

func (s *Service) InitQueue(settings *interfaces.Queue)  {
	if settings.Readable {
		s.InitConsumer(settings)
	}

	if settings.Writable {
		s.InitProducer(settings)
	}
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

func (s *Service) initProjectConfig()  {
	projectConfig := s.Project.LoadConfig()
	s.Config = projectConfig

}

func (s *Service) InitProducer(settings *interfaces.Queue) {
	name := settings.GetGroupName()
	color.Yellow(`init NSQ producer. 
	Group:%s 
	Topic: %s 
	Channel: %s`,  name, settings.Topic, settings.Channel,
	)
	if s.Config == nil {
		s.initProjectConfig()
	}
	producer, _ := nsq.StartProducer(nsq.ProducerConfig{
		Address: s.Config.NsqdNodes[0].IP,
	})


	if s.Producers == nil {
		s.Producers = Producers{
			name: &Producer{Worker: producer},
		}
	} else {
		s.Producers[name] = &Producer{Worker:producer}
	}

}

func (s *Service) InitConsumer(
	settings *interfaces.Queue,
	) *nsq.Consumer {

	name := settings.GetGroupName()
	color.Blue(`init NSQ consumer. 
	Group:%s 
	Topic: %s 
	Channel: %s`,  name, settings.Topic, settings.Channel,
	)

	if s.Config == nil {
		s.initProjectConfig()
	}
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:   settings.Topic,
		Channel: settings.Channel,
		Address: s.Config.NsqdNodes[0].IP,
		ReadTimeout: time.Duration(60) * time.Second,
		MaxInFlight: settings.Concurrent,
	})
	if s.Consumers == nil {
		s.Consumers = Consumers{
			name: {
				&Consumer{
					Topic:      settings.Topic,
					Channel:    settings.Channel,
					Concurrent: settings.Concurrent,
					Worker:     consumer,
				},
			},
		}
	} else {
		s.Consumers[name] = append(s.Consumers[name], &Consumer{
			Topic:      settings.Topic,
			Channel:    settings.Channel,
			Concurrent: settings.Concurrent,
			Worker:     consumer,
		})
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

func (s *Service) StopProducers()  {
	color.Yellow("Stop all producers ...")
	for producer := range s.Producers {
		s.Producers[producer].Worker.Stop()
	}

}

func (s *Service) Pub(
	name string,
	topic string,
	message string,
	) error {

	return s.Producers[name].Worker.PublishTo(topic, []byte(message))
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

func (s *Service) mergeConsumerChannels()[]<-chan interface{} {
	var outs [] <-chan interface{}

	for name := range s.Consumers {
		for _, consumer := range s.Consumers[name] {
			consumer.Name = name
			outs = append(outs, s.read(consumer))
		}
	}
	return outs

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


func  (s *Service) Collect(
	ctx context.Context,
	maxSize int,
	maxDuration time.Duration,
	in <-chan *nsq.Message,
	) <-chan []*nsq.Message {

		out := make(chan []*nsq.Message)
		go func() {
			for {
				is, open := s.collect(ctx, maxSize, maxDuration, in)
				if is != nil {
					out <- is
				}
				if !open {
					close(out)
					return
				}
			}
		}()
		return out
}

func (s *Service) collect(
	ctx context.Context,
	maxSize int,
	maxDuration time.Duration,
	in <-chan *nsq.Message,
	) ([]*nsq.Message, bool) {

		var buffer []*nsq.Message
		timeout := time.After(maxDuration)
		for {
			lenBuffer := len(buffer)
			select {
			case <-ctx.Done():
				bs, open := s.collect(context.Background(), maxSize, 100*time.Millisecond, in)
				return append(buffer, bs...), open
			case <-timeout:
				return buffer, true
			case i, open := <-in:
				if !open {
					return buffer, false
				} else if lenBuffer < maxSize-1 {
					// There is still room in the buffer
					buffer = append(buffer, i)
				} else {
					// There is no room left in the buffer
					return append(buffer, i), true
				}
			}
		}
	}

func (s *Service) batchRead(
	ctx context.Context,
	maxSize int,
	maxDuration time.Duration,
	in <-chan *nsq.Message,
	) (<-chan []*nsq.Message) {
		out := make(chan []*nsq.Message)
		is, _ := s.collect(ctx, maxSize, maxDuration, in)

		if is != nil {
			select {
			// Cancel all inputs during shutdown
			case <-ctx.Done():
				//processor.Cancel(is, ctx.Err())
			// Otherwise Process the inputs
			default:
				out <- is
			}
		}

		return out

	}




func (s *Service) BatchRead()  {

	type Args struct {
		maxSize     int
		maxDuration time.Duration
		in          []interface{}
		inDelay     time.Duration
		ctxTimeout  time.Duration
	}

	args := &Args{
		maxSize:     10,
		maxDuration: 100,
		in:          nil,
		ctxTimeout:  20,
	}

	// Create the context
	ctx, cancel := context.WithTimeout(context.Background(), args.ctxTimeout)
	defer cancel()

	const maxTestDuration = time.Second

	stream := pipeline.Merge(s.mergeConsumerChannels()...)

	// Collect responses
	collect := pipeline.Collect(ctx, args.maxSize, args.maxDuration, stream)
	timeout := time.After(maxTestDuration)
	var outs []interface{}
	var isOpen bool

	loop:
		for {
			select {
			case out, open := <-collect:
				if !open {
					isOpen = false
					break loop
				}
				isOpen = true
				outs = append(outs, out)
			case <-timeout:
				break loop
			}
		}
		if isOpen {

		}


	}