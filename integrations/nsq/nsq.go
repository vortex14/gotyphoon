package nsq

import (
	"context"
	"github.com/deliveryhero/pipeline"
	"github.com/fatih/color"
	"github.com/segmentio/nsq-go"
	"github.com/vortex14/gotyphoon/config"
	"github.com/vortex14/gotyphoon/interfaces"
	"strings"
	"time"
)

type Consumer struct {
	Name string
	Topic string
	Channel string
	Concurrent int
	Worker *nsq.Consumer
}

type Producer struct {
	Worker *nsq.Producer
}

type Producers map[string] *Producer
type Consumers map[string] [] *Consumer

type Yield struct {
	Name string
	Topic string
	Channel string
	Msg *nsq.Message
}

type Service struct {
	Config *config.Config
	Project interfaces.Project
	Producers Producers
	Consumers Consumers
}

func (s *Service) TestConnect() bool  {
	status := false
	var projectConfig config.Config
	if s.Config != nil {
		projectConfig = *s.Config
	} else {
		projectConfig = s.Project.LoadConfig().Config
	}

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

func (s *Service) initProjectConfig()  {
	projectConfig := s.Project.LoadConfig().Config
	s.Config = &projectConfig

}

func (s *Service) InitProducer(name string) {
	color.Yellow("init %s nsq producer", name)
	if s.Config == nil {
		s.initProjectConfig()
	}
	producer, _ := nsq.StartProducer(nsq.ProducerConfig{
		Address: s.Config.NsqdNodes[0].IP,
	})

	s.Producers = Producers{
		name: &Producer{Worker: producer},
	}

}

func (s *Service) InitConsumer(
	name string,
	topic string,
	channel string,
	concurrent int,
	) *nsq.Consumer {

	color.Yellow("init %s nsq consumer", name)
	if s.Config == nil {
		s.initProjectConfig()
	}
	consumer, _ := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:   topic,
		Channel: channel,
		Address: s.Config.NsqdNodes[0].IP,
		ReadTimeout: time.Duration(60) * time.Second,
		MaxInFlight: concurrent,
	})
	if s.Consumers == nil {
		s.Consumers = Consumers{
			name: {
				&Consumer{
					Topic:      topic,
					Channel:    channel,
					Concurrent: concurrent,
					Worker:     consumer,
				},
			},
		}
	} else {
		s.Consumers[name] = append(s.Consumers[name], &Consumer{
			Topic:      topic,
			Channel:    channel,
			Concurrent: concurrent,
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
			out <- &Yield{
				Msg:     &msg,
				Name: consumer.Name,
				Topic:   consumer.Topic,
				Channel: consumer.Channel,
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

func (s *Service) mergeRead(ins ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	for i := range ins {
		go func(in <-chan interface{}) {
			for i := range in {
				out <- i
			}
		}(ins[i])
	}
	return out
}


func (s *Service) Read() {


	for source := range pipeline.Merge(s.mergeConsumerChannels()...) {
		yield := source.(*Yield)

		color.Red("%+v", yield)
		yield.Msg.Finish()

	}

	//for yield := range s.mergeRead(s.mergeConsumerChannels()...) {
	//	color.Green(`
	//
	//	Topic read: %s
	//	Channel read: %s
	//	Body: %s
	//	Name: %s
	//
	//	`,
	//	yield.Topic,
	//	yield.Channel,
	//	string(yield.Msg.Body),
	//	yield.Name,
	//	)
	//
	//	yield.Msg.Finish()
	//}

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



//	//ctx, cancel := context.WithTimeout(context.Background(), args.ctxTimeout)
//
//
//	for yield := range s.mergeRead(s.mergeConsumerChannels()...) {
//		color.Green(`
//
//		Topic read: %s
//		Channel read: %s
//		Body: %s
//		Name: %s
//
//		`,
//			yield.Topic,
//			yield.Channel,
//			string(yield.Msg.Body),
//			yield.Name,
//		)
//
//		yield.Msg.Finish()
//	} }
