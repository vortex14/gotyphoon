package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/utils/code"
	"os"

	"github.com/vortex14/gotyphoon"
	"github.com/vortex14/gotyphoon/elements/models/awaitable"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/integrations/nsq"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/utils"
)

func init()  {
	log.InitD()
}

const (
	TASKS   = "tasks"
	TOPIC   = "agent"
	CHANNEL = "tasks"
)

type Task struct {
	Image string `json:"image"`
	Tag   string `json:"tag"`
}

type Agent struct {
	singleton.Singleton
	awaitable.Object

	ConfigFile string
	project    *typhoon.Project
	LOG        interfaces.LoggerInterface

	messageBus *nsq.Service
	Topic      string
	Channel    string
}

func (a *Agent) GetTasks()  {
	a.LOG.Debug("get tasks")
	config := a.project.LoadConfig()
	setting := &interfaces.Queue{Topic: config.Topic, Concurrent: config.Concurrent, Channel: config.Channel}
	setting.SetGroupName(TASKS)

	consumer := a.messageBus.InitConsumer(setting)

	var count int
	count = 0
	for msg := range consumer.Messages() {

		color.Yellow("%s", msg.Body)
		var task Task
		count += 1
		_ = json.Unmarshal(msg.Body, &task)

		a.LOG.Debug(fmt.Sprintf("№-%d %+v", count, task))

		msg.Finish()
	}
}

func (a *Agent) Init()  {
	a.Construct(func() {
		a.LOG = log.New(log.D{"agent": "ci-agent"})
		a.project = &typhoon.Project{
			ConfigFile: a.ConfigFile,
			Path:       utils.GetCurrentDir(),
		}

		nsqService := &nsq.Service{Project: a.project}
		status := nsqService.Ping()
		if !status { color.Red("Connection failed to NSQ"); os.Exit(1) } else {
			color.Green("NSQ connected !")
		}
		a.messageBus = nsqService

		a.Add()
		go a.GetTasks()
	})
}

func CommentedCode(marker string, path string)  {

}

func main()  {
	println("start agent ")

	matchCode := "ignore for building amd64-linux"

	//fmt.Println(path)
	startDir := "../../../"
	excludeDirs := map[string]bool{"vendor": true, ".git": true, "tmp": true, ".idea": true}
	code.UncommentDir(startDir, matchCode, excludeDirs)
	//utils.CommentDir(startDir, matchCode, excludeDirs)
	return

	agent := Agent{
		ConfigFile: "config.agent.local.yaml",
	}
	agent.Init()


	agent.Await()

}