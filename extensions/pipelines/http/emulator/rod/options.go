package rod

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/proto"
)

type EventOptions struct {
	NetworkResponseReceived bool
	Page                    *rod.Page
}

func (e *EventOptions) Wait() {

	if e.NetworkResponseReceived {
		er := proto.NetworkResponseReceived{}
		wait := e.Page.WaitEvent(&er)
		wait()
	}

}

type Options struct {
	Debug       bool
	Proxy       string
	Timeout     time.Duration
	Device      devices.Device
	RandomAgent bool
}

type DetailsOptions struct {
	Options

	EventOptions  EventOptions
	ProxyRequired bool
	Click         bool
	MustElement   string
	Input         string
	SleepAfter    float32
}
