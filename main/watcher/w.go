package main

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	
	BaseWatcher "github.com/vortex14/gotyphoon/elements/models/watcher"
	Watcher "github.com/vortex14/gotyphoon/extensions/models/watcher"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)

func main() {
	path := utils.GetCurrentDir()
	pathJoin := filepath.Join(path, "main", "watcher")
	w := Watcher.Watcher{
		Watcher: BaseWatcher.Watcher{
			Path: pathJoin,
		},
		Callback: func(log interfaces.LoggerInterface, event *fsnotify.Event) {
			log.Info(event)
		},
	}
	w.Watch()
	println(pathJoin)
}
