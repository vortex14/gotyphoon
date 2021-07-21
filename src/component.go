package typhoon

import (
	"bytes"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-cmd/cmd"
	"github.com/go-logfmt/logfmt"
	"github.com/vortex14/gotyphoon/src/interfaces"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Directory struct {
	Path string
}

type Worker struct {
	Command 	  string
	Args    	  []string
	Cmd			  *cmd.Cmd
	Status 		  chan bool
}

type Component struct {
	Path string
	Name string
	isException bool
	isDebug bool
	Active bool
	Worker *Worker
	Promise sync.WaitGroup
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
	color.Yellow("initStart %s", c.Name)

	pathExecute := fmt.Sprintf("%s.py", c.Name)
	configArg := fmt.Sprintf("--config=%s", project.GetConfigFile())
	logLevelArg := fmt.Sprintf("--level=%s", project.GetLogLevel())
	projectNameArg := fmt.Sprintf("--project_name=%s", project.GetName())
	c.Worker = &Worker{Command: "python3.8", Args: []string{pathExecute, configArg, logLevelArg, projectNameArg }}
	//c.Path = fmt.Sprintf("%s/project/%s", project.GetProjectPath(), c.Name )
	c.Worker.Run(project)

	c.Worker.Cmd.Start()
	c.Worker.Cmd.Status()
	c.Active = true
	//c.AddPromise()
	go c.Logging()
}


func (c *Component) Close(project interfaces.Project)  {
	defer project.PromiseDone()
	c.Stop(project)
}

func exec_cmd(cmd *exec.Cmd) {
	var waitStatus syscall.WaitStatus
	err := cmd.Run()

	if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		waitStatus = exitError.Sys().(syscall.WaitStatus)
		fmt.Printf("Error during killing (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		fmt.Printf("Port successfully killed (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
}


func (c *Component) Stop(project interfaces.Project)  {
	status := c.Worker.Cmd.Status()
	color.Green("%s status.PID %s", status.PID, c.Name)
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

	command := fmt.Sprintf("lsof -i :%s", port)

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

func (c Component) Restart(project *Project)  {
	color.Red("Restart component %s ...", c.Name)
	c.Stop(project)
	c.Start(project)

	project.components.ActiveComponents[c.Name] = &c
}

func (c *Component) GetName() string {
	return c.Name
}



func (d *Directory) GetDataFromDirectory(path string) interfaces.MapFileObjects {
	currentData := make(interfaces.MapFileObjects, 0)


	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		typeFile := "file"
		if file.IsDir() {
			typeFile = "dir"
		}

		currentData[file.Name()] = &interfaces.FileObject{

			Type: typeFile,
			Path: file.Name(),
		}

	}



	return currentData
}

func (c *Component) CheckDirectory(required []string, pathComponent string) bool  {
	var status = true

	color.Yellow("checking %s", pathComponent)
	dir := &Directory{
		Path: pathComponent,
	}
	dataDir := dir.GetDataFromDirectory(dir.Path)
	for _, reqFile := range required {
		if _, ok := dataDir[reqFile]; !ok {
			color.Red("%s not exist in %s", reqFile, pathComponent)
			status = false
		} else {
			color.Green("checked %s in %s", reqFile, pathComponent)
		}


	}


	return status
}

func (d *Directory) IsExistDir(path string) bool  {
	var status = false
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		status = true
	}

	return status
}

func (c *Component) CheckComponent() bool {
	var status = false

	pathComponent := fmt.Sprintf("project/%s",c.Name)



	if _, err := os.Stat(pathComponent); !os.IsNotExist(err) {

		if c.Name == "fetcher" {
			required := []string{"executions", "responses", "__init__.py"}

			status = c.CheckDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Fetcher dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if c.Name == "processor" {

			required := []string{"executable", "__init__.py"}
			status = c.CheckDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Processor dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if c.Name == "scheduler" {
			required := []string{"__init__.py"}
			status = c.CheckDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Scheduler dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if c.Name == "donor" {
			required := []string{"__init__.py", "v1", "routes.py"}
			status = c.CheckDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Scheduler dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}

		} else if c.Name == "result_transporter" {
			required := []string{"__init__.py", "consumers"}
			status = c.CheckDirectory(required, pathComponent)
			logVal := fmt.Sprintf("Scheduler dir is %t", status)
			if status == true {
				color.Green(logVal)
			} else {
				color.Red(logVal)
			}
		}
	} else {
		color.Red("path %s doesn't exist", c.Name)
	}

	fileName := fmt.Sprintf("%s.py", c.Name)
	required := []string{fileName}
	status = c.CheckDirectory(required, ".")
	logVal := fmt.Sprintf("%s.py is %t", c.Name, status)

	if status == true {
		color.Green(logVal)
	} else {
		color.Red(logVal)
	}


	return status
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
				fmt.Printf(`
-DEBUG-----------
`)
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


			if strings.Contains(line, ">>>!") {
				c.isException = true
				color.Red("Component %s has error.", c.Name)
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



					if c.Name == "processor" {
						color.Yellow("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if c.Name == "result_transporter" {
						color.Green("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if c.Name == "fetcher" {
						color.Blue("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if c.Name == "donor" {
						color.HiBlackString("%s = %s", logDataMap.Key(), logDataMap.Value())
					} else if c.Name == "scheduler" {
						color.Cyan("%s = %s", logDataMap.Key(), logDataMap.Value())
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
			io.Copy(os.Stderr, bytes.NewBufferString(errLog))
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

			if ok != true || status == false {

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