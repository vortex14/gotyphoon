package nsq

import (
	"context"
	"time"

	"github.com/segmentio/nsq-go"
)

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

