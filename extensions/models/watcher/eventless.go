package watcher

import (
	"github.com/vortex14/gotyphoon/log"
	"time"

	lessWatcher "github.com/radovskyb/watcher"
	"github.com/vortex14/gotyphoon/elements/models/watcher"
	"github.com/vortex14/gotyphoon/interfaces"
)

type EventFSLess struct {
	watcher.Watcher
	Callback func(log interfaces.LoggerInterface, event *lessWatcher.Event)
	Timeout  int32
}

func (w *EventFSLess) Watch() {
	w.Construct(func() {
		LOG := log.New(log.D{"watch": w.Path})
		ww := lessWatcher.New()
		// SetMaxEvents to 1 to allow at most 1 event's to be received
		// on the Event channel per watching cycle.
		//
		// If SetMaxEvents is not set, the default is to send all events.
		ww.SetMaxEvents(1)

		// Only notify for
		ww.FilterOps(
			lessWatcher.Rename,
			lessWatcher.Move,
			lessWatcher.Write,
			lessWatcher.Remove,
		)

		go func() {
			for {
				select {
				case event := <-ww.Event:
					// pass to callback
					w.Callback(LOG, &event)
				case err := <-ww.Error:
					LOG.Error(err)
				case <-ww.Closed:
					return
				}
			}
		}()

		// Watch path recursively for changes.
		if err := ww.AddRecursive(w.Path); err != nil {
			LOG.Error(err)
		}

		// Start the watching process - it'll check for changes every ±100ms.
		if err := ww.Start(time.Millisecond * time.Duration(w.Timeout)); err != nil {
			LOG.Error(err)
		}

	})
}
