package nsq

import (
	"context"
	"github.com/segmentio/nsq-go"
	"time"

	"github.com/deliveryhero/pipeline"
)

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

func (s *Service) batchRead(
	ctx context.Context,
	maxSize int,
	maxDuration time.Duration,
	in <-chan *nsq.Message,
) <-chan []*nsq.Message {
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

