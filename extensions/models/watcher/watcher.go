package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/vortex14/gotyphoon/elements/models/watcher"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type Watcher struct {
	watcher.Watcher
	Callback func(log interfaces.LoggerInterface, event *fsnotify.Event)
}

func (w *Watcher) Watch() {
	LOG := log.New(log.D{"watch": w.Path})
	watcher, _ := fsnotify.NewWatcher()
	err := watcher.Add(w.Path)
	if err != nil {
		LOG.Error(err.Error())
		return
	}
	LOG.Debug("start watching ...")
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				LOG.Debug(event.Name)
				if w.IsIgnore(event.Name) {
					continue
				}
				if w.Callback != nil {
					w.Callback(LOG, &event)
				}

			case err := <-watcher.Errors:
				LOG.Error("ERROR---->", err.Error())
			}
		}
	}()

	<-done
}
