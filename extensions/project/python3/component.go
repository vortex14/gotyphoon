package python3

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	TyLog "github.com/vortex14/gotyphoon/log"

	"github.com/fatih/color"
	"github.com/go-logfmt/logfmt"

	. "github.com/vortex14/gotyphoon/extensions/models/cmd"
	"github.com/vortex14/gotyphoon/interfaces"

	"github.com/vortex14/gotyphoon/elements/forms"
)

type Component struct {
	forms.Component

	CMDProcess *Command
}

func (c *Component) Init(project interfaces.Project) {
	projectNameArg := fmt.Sprintf("--project_name=%s", project.GetName())
	logLevelArg := fmt.Sprintf("--level=%s", project.GetLogLevel())
	configArg := fmt.Sprintf("--config=%s", project.GetConfigPath())

	cmd := &Command{
		Cmd:  "python3.8",
		Dir:  project.GetProjectPath(),
		Args: []string{c.FileExt, projectNameArg, configArg, logLevelArg},
	}

	err := cmd.Run()
	c.Active = true
	c.CMDProcess = cmd
	if err != nil {
		color.Red(err.Error())
		return
	}

	c.Add()
	go c.Logging()
}

func (c *Component) Start(project interfaces.Project) {
	c.Construct(func() {
		c.LOG = TyLog.New(TyLog.D{"component": c.Name})
		defer project.PromiseDone()

		c.Init(project)
		//executable := filepath.Join(project.GetProjectPath(), c.FileExt)

	})
}

func (c *Component) Close(project interfaces.Project) {
	c.LOG.Info("Close project")
	defer project.PromiseDone()

	c.CMDProcess.Close()
	c.Done()
	c.Await()
	c.CMDProcess.Await()
	c.Active = false
	c.LOG.Debug("Close")

}

func (c *Component) Restart(project interfaces.Project) {
	//defer project.PromiseDone()
	color.Red("Restart component %s ...", c.Name)
	project.AddPromise()
	c.Close(project)
	c.Init(project)
	c.Start(project)

	color.Green("restarted !")
}

func (c *Component) CheckComponent() bool {
	status := false
	var (
		logVal        string
		componentName string
		required      []string
	)

	switch c.Name {
	case interfaces.FETCHER:
		componentName = interfaces.TYPHOON2PYTHON2FETCHER
		required = []string{"executions", "responses", "__init__.py"}
		logVal = "Check Fetcher dir"
	case interfaces.PROCESSOR:
		componentName = interfaces.TYPHOON2PYTHON2PROCESSOR
		required = []string{"executable", "__init__.py"}
		logVal = "Check processor dir"
	case interfaces.DONOR:
		componentName = interfaces.TYPHOON2PYTHON2DONOR
		required = []string{"__init__.py", "v1", "routes.py"}
		logVal = "Check donor dir"
	case interfaces.TRANSPORTER:
		componentName = interfaces.TYPHOON2PYTHON2TRANSPORTER
		required = []string{"__init__.py", "consumers"}
		logVal = "Check transporter dir"
	case interfaces.SCHEDULER:
		componentName = interfaces.TYPHOON2PYTHON2SCHEDULER
		required = []string{"__init__.py"}
		logVal = "Check scheduler dir"
	default:
		color.Red("Component not found %s", c.Name)
		os.Exit(1)
	}

	c.InitFolder(componentName)

	if !c.Folder.IsExist(".") {
		color.Red("Component: %s Path %s doesn't exist", c.Name, c.Folder.Path)
		return false
	}

	_, status = c.Folder.IsExists(required)

	color.Yellow("Path: %s Component: %s Status: %t", c.Folder.Path, c.Name, status)

	if status {
		color.Green(logVal)
	} else {
		color.Red(logVal)
	}

	fileNameExt := fmt.Sprintf("%s.py", componentName)
	color.Yellow("Check file %s", fileNameExt)
	required = []string{fileNameExt}
	_, status = c.ProjectFolder.IsExists(required)
	logVal = fmt.Sprintf("%s.py is %t", componentName, status)

	if status {
		color.Green(logVal)
		c.FileExt = fileNameExt
	} else {
		color.Red(logVal)
	}

	return status
}

func (c *Component) Logging() {
	Info := color.New(color.FgWhite, color.BgBlack, color.Bold).SprintFunc()
	c.IsException = false
	c.IsDebug = false
	for {
		select {
		case line, open := <-c.CMDProcess.Output:
			if !open {
				continue
			}

			if strings.Contains(line, "@debug") {
				fmt.Printf(TyLog.OWL)
				c.IsDebug = true
				continue
			}
			if strings.Contains(line, "/debug") {
				c.IsDebug = false
				fmt.Printf(`
-/DEBUG-----------


`)
				continue
			}
			if c.IsDebug {
				color.Green(line)
				continue
			}

			if strings.Contains(line, ">>>!") || strings.Contains(line, "level=ERROR") && !c.IsException {
				c.IsException = true
				fmt.Printf(TyLog.DINOSAUR, c.Name)
				color.Red(line)
				continue
			}

			if strings.Contains(line, "!<<<") {
				c.IsException = false
				color.Red(line)
				fmt.Printf(`
------------
`)
				continue
			}

			if c.IsException {
				color.Red(line)

				continue
			}

			color.Cyan(line)
			fmt.Printf(`%s Logs ...
`, Info(c.Name))
			logDataMap := logfmt.NewDecoder(strings.NewReader(line))

			for logDataMap.ScanRecord() {
				for logDataMap.ScanKeyval() {

					switch c.Name {
					case interfaces.FETCHER:
						color.Blue("%s = %s", logDataMap.Key(), logDataMap.Value())
					case interfaces.PROCESSOR:
						color.Yellow("%s = %s", logDataMap.Key(), logDataMap.Value())
					case interfaces.SCHEDULER:
						color.Cyan("%s = %s", logDataMap.Key(), logDataMap.Value())
					case interfaces.TRANSPORTER:
						color.Green("%s = %s", logDataMap.Key(), logDataMap.Value())
					case interfaces.DONOR:
						color.Magenta("%s = %s", logDataMap.Key(), logDataMap.Value())
					}

				}

			}
			if logDataMap.Err() != nil {
				//color.Red("Invalid Log format. Don't use = . Broken line: %s",line)
				//panic(d.Err())
				continue
			}
			fmt.Printf(`
------------
`)
			continue
		case line, open := <-c.CMDProcess.OutputErr:
			if !open {
				continue
			}
			errLog := ""
			_, err := io.Copy(os.Stderr, bytes.NewBufferString(errLog))
			if err != nil {
				color.Red("%s", err.Error())
				return
			}
			color.Red(" %s error: %s", c.Name, line)

		}

	}

}
