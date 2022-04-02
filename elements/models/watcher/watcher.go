package watcher

import (
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"strings"
)

type Watcher struct {
	singleton.Singleton
	Path         string
	ExcludeMatch []string
}

func (w *Watcher) IsIgnore(eventName string) bool {
	status := false

	for _, excTerm := range w.ExcludeMatch {
		if strings.Contains(eventName, excTerm) {
			status = true
			break
		}
	}

	return status
}

func (w *Watcher) Watch() {

}
