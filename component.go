package typhoon

import (
	"bytes"
	"context"
	"fmt"
	"github.com/vortex14/gotyphoon/elements/models/folder"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	TyLog "github.com/vortex14/gotyphoon/log"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	//"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/go-cmd/cmd"
	"github.com/go-logfmt/logfmt"

	. "github.com/vortex14/gotyphoon/extensions/models/cmd"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Worker struct {
	Command 	  string
	Args    	  []string
	Cmd			  *cmd.Cmd
	Status 		  chan bool
}

type file struct {
	FileExt string
	Language string
}

type Component struct {
	singleton.Singleton
	folder *folder.Folder

	ProjectPath string

	file
	Name string
	Active bool
	isDebug bool
	Worker *Worker
	isException bool
	Promise sync.WaitGroup
	Producers interfaces.Producers
	Pipelines interfaces.Pipelines
	Consumers interfaces.Consumers
	QueuesSettings interfaces.Queue
}

func (c *Component) AddPromise()  {
	c.Promise.Add(1)
}

func (c *Component) PromiseDone()  {
	c.Promise.Done()
}

func (c *Component) WaitPromises()  {
	c.Promise.Wait()
}

func (c *Component) Start(project interfaces.Project)  {

	//executable := filepath.Join(project.GetProjectPath(), c.FileExt)

	projectNameArg := fmt.Sprintf("--project_name=%s", project.GetName())
	logLevelArg := fmt.Sprintf("--level=%s", project.GetLogLevel())
	configArg := fmt.Sprintf("--config=%s", project.GetConfigPath())

	
	cmd := Command{
		Cmd:       "python3.8",
		Dir:       project.GetProjectPath(),
		Args:      []string{c.FileExt, projectNameArg, configArg, logLevelArg},
	}

	err := cmd.Run()
	if err != nil {
		color.Red(err.Error())
		return 
	}

	go func() {
		for line := range cmd.Output {
			println(line)
		}

	}()

	cmd.Await()

	//c.Worker = &Worker{Command: "python3.8", Args: []string{executable, configArg, logLevelArg, projectNameArg }}
	////c.Path = fmt.Sprintf("%s/project/%s", project.GetProjectPath(), c.Name )
	//c.Worker.Run(project)
	//
	//c.Worker.Cmd.Start()
	//c.Worker.Cmd.Status()
	//c.Active = true
	////c.AddPromise()
	//go c.Logging()
}


func (c *Component) Close(project interfaces.Project)  {
	defer project.PromiseDone()
	c.Stop(project)
}
//
//func exec_cmd(cmd *exec.Cmd) {
//	var waitStatus syscall.WaitStatus
//	err := cmd.Run()
//
//	if err != nil {
//			os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
//	}
//	if exitError, ok := err.(*exec.ExitError); ok {
//		waitStatus = exitError.Sys().(syscall.WaitStatus)
//		fmt.Printf("Error during killing (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
//	} else {
//		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
//		fmt.Printf("Port successfully killed (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
//	}
//}


func (c *Component) Stop(project interfaces.Project)  {
	status := c.Worker.Cmd.Status()
	color.Green("%d status.PID %s", status.PID, c.Name)
	//if !IsClosed(c.Worker.Status){
	c.Worker.Status <- false
	//}
	c.Active = false
	port := project.GetComponentPort(c.Name)

	//
	//if runtime.GOOS == "windows" {
	//	command := fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s).OwningProcess -Force", port)
	//	exec_cmd(exec.Command("Stop-Process", "-Id", command))
	//} else {
	//	command := fmt.Sprintf("lsof -i :%s", port)
	//	exec_cmd(exec.Command("bash", "-c", command))
	//}

	command := fmt.Sprintf("lsof -i :%d", port)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "bash", "-c", command).Run(); err != nil {
		color.Green("commands done !")
		// This will fail after 100 milliseconds. The 5 second sleep
		// will be interrupted.
	}



	//cmdSource := exec.Command("bash", "-c", command)
	//var out bytes.Buffer
	//cmdSource.Stdout = &out
	//_ = cmdSource.Run()
	//
	//err := cmdSource.Wait()
	//if err != nil {
	//	color.Green("Yes: %s", err)
	//}
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("in all caps: %q\n", out.String())

	//err := cmdSource.Start()
	//if err != nil {
	//	color.Red("%s", err)
	//}
	//
	//data, errs := cmdSource.Output()
	//color.Red("test err %s", errs)
	//color.Green("Port: %s flushed", string(rune(port)))

	//
	//_, err := cmdSource.CombinedOutput()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//s := bufio.NewScanner(bytes.NewReader(bs))
	//
	//for s.Scan() {
	//	kv := strings.SplitN(s.Text(), "=", 2)
	//	if strings.Contains(strings.ToLower(kv[0]), "typhoon") {
	//		os.Setenv(kv[0], kv[1])
	//	}
	//}

	errKill := syscall.Kill(-status.PID, syscall.SIGKILL)
	if errKill == nil {
		color.Green("%s killed", c.Name)
		//color.Red("Error kill :%s, component: %s", errKill, c.Name)
	}





	color.Red("component %s was be closed", c.Name)

}

func (c *Component) Restart(project *Project)  {
	color.Red("Restart component %s ...", c.Name)
	c.Stop(project)
	c.Start(project)

	project.components.ActiveComponents[c.Name] = c
}

func (c *Component) GetName() string {
	return c.Name
}

func (c *Component) initFolder()  {


	c.folder = &folder.Folder{Path: c.ProjectPath}
}

func (c *Component) CheckComponent() bool {
	status := false
	var (
		logVal string
		componentName string
		required []string
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

	c.folder = &folder.Folder{Path: filepath.Join(c.ProjectPath, "project", componentName)}

	if !c.folder.IsExist(".") {
		color.Red("Component: %s Path %s doesn't exist", c.Name, c.folder.Path)
		return false
	}



	_, status = c.folder.IsExists(required)

	color.Yellow("Path: %s Component: %s Status: %t", c.folder.Path, c.Name, status)


	if status {
		color.Green(logVal)
	} else {
		color.Red(logVal)
	}

	fileNameExt := fmt.Sprintf("%s.py", componentName)
	color.Yellow("Check file %s", fileNameExt)
	required = []string{fileNameExt}
	_, status = c.folder.IsExists(required)
	logVal = fmt.Sprintf("%s.py is %t", componentName, status)

	if status {
		color.Green(logVal)
		c.FileExt = fileNameExt
	} else {
		color.Red(logVal)
	}


	return status
}

func (c *Component) InitConsumers(project interfaces.Project)  {
	config := project.LoadConfig()
	queueSettings :=  config.TyComponents.Fetcher.Queues
	color.Yellow("current fetcher settings %+v", queueSettings)
	color.Yellow("InitConsumers for %s", c.Name)
}

func (c *Component) InitProducers()  {

}


func (c *Component) StopConsumers()  {

}

func (c *Component) StopProducers()  {

}

func (c *Component) RunQueues() {

}

func (w *Worker) Run(project interfaces.Project) {
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	envCmd := cmd.NewCmdOptions(cmdOptions, w.Command, w.Args...)
	w.Cmd = envCmd
	w.Status = make(chan bool, 1)
	w.Status <- true
	projectEnv := fmt.Sprintf("PYTHONPATH=%s:%s", project.GetEnvSettings(), project.GetProjectPath())
	newEnv := append(os.Environ(), projectEnv)
	envCmd.Env = newEnv
}

func IsClosed(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func (c *Component) Logging()  {
	Info := color.New(color.FgWhite, color.BgBlack, color.Bold).SprintFunc()
	c.isException = false
	c.isDebug = false
	for {
		select {
		case line, open := <-c.Worker.Cmd.Stdout:
			if !open {
				continue
			}

			if strings.Contains(line, "@debug") {
				fmt.Printf(TyLog.OWL)
				c.isDebug = true
				continue
			}
			if strings.Contains(line, "/debug") {
				c.isDebug = false
				fmt.Printf(`
-/DEBUG-----------


`)
				continue
			}
			if c.isDebug {
				color.Green(line)
				continue
			}


			if strings.Contains(line, ">>>!") || strings.Contains(line, "level=ERROR") && !c.isException {
				c.isException = true
				fmt.Printf(TyLog.DINOSAUR, c.Name)
				color.Red(line)
				continue
			}

			if strings.Contains(line, "!<<<") {
				c.isException = false
				color.Red(line)
				fmt.Printf(`
------------
`)
				continue
			}

			if c.isException {
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
		case line, open := <-c.Worker.Cmd.Stderr:
			if !open {
				continue
			}
			errLog := ""
			_, err := io.Copy(os.Stderr, bytes.NewBufferString(errLog))
			if err != nil {
				color.Red("%s", err.Error())
				return
			}
			//errLog = fmt.Sprintf("Component: %s; %s , %s", w.Name, errLog, line)
			//color.Red(errLog)
			color.Red(" %s error: %s",c.Name, line)
			//err := c.Worker.Cmd.Stop()

			//if err != nil {
			//	color.Red(" %s error: %s",c.Name, line)
				//fmt.Fprintln(os.Stderr, line)
			//}
			//close(w.Status)

			//color.Red("Return from Logging. Component: %s", w.Name)
			//status := w.Cmd.Status()
			//errKill := syscall.Kill(-status.PID, syscall.SIGKILL)
			//if errKill != nil {
			//	color.Red("Error kill :%s, component: %s", errKill, w.Name)
			//}
			continue
		case status, ok := <-c.Worker.Status:

			if !ok || !status {

				err := c.Worker.Cmd.Stop()

				if err != nil {
					color.Red("Component: %s ,Err: %s", c.Name, err)
				}

				//
				//
				//if !IsClosed(c.Worker.Status) {
				//	close(c.Worker.Status)
				//}

				//c.Promise.Done()

				color.Blue("%s logging done ... ", c.Name)

				return

			}
			continue

		}

	}


}