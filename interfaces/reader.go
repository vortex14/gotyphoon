package interfaces

import "github.com/segmentio/nsq-go"

type YieldNsqMessage struct {
	*Yield
	Msg *nsq.Message
}

