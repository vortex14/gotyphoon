package main

import (
	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/elements/forms"
	"os"
	"path/filepath"
	"strings"

	pyTyphoon "github.com/vortex14/gotyphoon/extensions/project/python3"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/utils"

	. "github.com/vortex14/gotyphoon/elements/models/osignal"
)

func init() {
	log.InitD()
}

func getProjectPath() string {
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

func main() {
	color.Red("Project path >>>>>>>>> %s <<<<<<<<<", getProjectPath())
	project := (&pyTyphoon.Project{
		Project: forms.Project{
			Name:              "test-project-cmd",
			Path:              getProjectPath(),
			SelectedComponent: []string{interfaces.PROCESSOR, interfaces.FETCHER},
			ConfigFile:        "config.local.yaml",
			AutoReload:        true,
		},
	}).Run()

	(&OSignal{
		Callback: func(logger interfaces.LoggerInterface, sig os.Signal) {
			project.Close()
			project.WaitPromises()
			logger.Warning("close !")
		},
	}).Wait()

}
