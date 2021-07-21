package typhoon

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/vortex14/gotyphoon/src/interfaces"
	"github.com/vortex14/gotyphoon/src/utils"
	"os"
	"path/filepath"
)




var watcher *fsnotify.Watcher

func watchDir(path string, fi os.FileInfo, err error) error {

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func CreateProject() {

}

func WatchTest()  {
	color.Green("watch for project ..")
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	// starting at the root of the project, walk each file/directory searching for
	// directories

	if err := filepath.Walk("project", watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done

}



func ParseLogData(fileObject *interfaces.FileObject) error {
	u := utils.Utils{}
	err := u.ParseLog(fileObject)


	if err != nil {

		color.Red("Error %s", err)
		os.Exit(0)

	}

	return nil


}
