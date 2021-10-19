package main

import (
	"os"
	"path/filepath"
	"strings"

	typhoon "github.com/vortex14/gotyphoon"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/utils"

	. "github.com/vortex14/gotyphoon/elements/models/osignal"
)

func init()  {
	log.InitD()
}

func getProjectPath()  string {
	herePath := utils.GetCurrentDir()
	var projectPath string
	// for __main__
	if strings.Contains(herePath, "python3") {
		projectPath = filepath.Join(herePath, "tp")
	} else {
		// for .air
		projectPath = filepath.Join(herePath, "main", "project", "python3", "tp")
	}
	return projectPath
}

func main()  {
	project := (&typhoon.Project{
		Path: getProjectPath(),
		SelectedComponent: []string{interfaces.PROCESSOR},
		ConfigFile:        "config.local.yaml",
		AutoReload:        true,
	}).Run()

	(&OSignal{
		Callback: func(logger interfaces.LoggerInterface, sig os.Signal) {
			logger.Warning("close !")
			project.Close()
		},
	}).Wait()

}
