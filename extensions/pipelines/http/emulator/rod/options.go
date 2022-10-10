package rod

import (
	"github.com/go-rod/rod"
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
