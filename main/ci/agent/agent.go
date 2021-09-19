package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

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

func UnCommentCode(marker string, code string) string {
	lines := strings.Split(code, "\n")
	var unCommentedLineList []string
	for _, line := range lines {
		if strings.Contains(line, marker) { unCommentedLineList = append(unCommentedLineList, line);continue }
		if strings.Contains(line, "*/") { unCommentedLineList = append(unCommentedLineList, line); continue}
		uncommentLine := strings.Replace(line, "//", "", 1)
		unCommentedLineList = append(unCommentedLineList, uncommentLine)

	}
	unCommentedLines := strings.Join(unCommentedLineList, "\n")
	return unCommentedLines
}

func CommentCode(marker string, code string) string {
	lines := strings.Split(code, "\n")
	//var stopMarker bool
	var CommentedLineList []string

	isComment := true

	for _, line := range lines {
		if utils.IsStrContain(line, "package") { CommentedLineList = append(CommentedLineList, line); continue }


		if strings.Contains(line, " */") {
			isComment = false
			CommentedLineList = append(CommentedLineList, line)
		} else if strings. Contains(line, marker) {
			isComment = true
			CommentedLineList = append(CommentedLineList, line)
		} else if strings.Contains(line, fmt.Sprintf("// %s", marker)) {
			isComment = true
			CommentedLineList = append(CommentedLineList, line)
		} else if strings.Contains(line, "//") {
			CommentedLineList = append(CommentedLineList, line)

		} else if isComment {
			commentLine := fmt.Sprintf("//%s", line)
			CommentedLineList = append(CommentedLineList, commentLine)
		} else {
			CommentedLineList = append(CommentedLineList, line)
		}

		//else if strings.Contains(line, marker) {
		//	isComment = true
		//	CommentedLineList = append(CommentedLineList, line)
		//} else if strings.Contains(line, "*/") {
		//	isComment = false
		//	CommentedLineList = append(CommentedLineList, line)
		//} else {
		//	isComment = false
		//}



		//if firstLine || !stopMarker && isComment {
		//	commentLine := fmt.Sprintf("//%s", line)
		//	CommentedLineList = append(CommentedLineList, commentLine)
		//} else if !firstLine {
		//	if strings.Contains(line, marker) { stopMarker = false; isComment = true }
		//	CommentedLineList = append(CommentedLineList, line)
		//} else if strings.Contains(line, marker) { stopMarker = false; isComment = true; continue }
		//println(line)
		//firstLine = false




	}
	CommentedLines := strings.Join(CommentedLineList, "\n")
	return CommentedLines
}

func UncommentDir(startDir string, matchCode string, excludeDirs map[string]bool)  {
	_ = filepath.Walk(startDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil { return err }
			firstDir := utils.GetFirstDir(path)
			if _, ok := excludeDirs[firstDir]; ok { return nil }
			if info.IsDir() { return nil}
			contentFileCode := utils.ReadFile(path)
			marker := fmt.Sprintf("/* %s", matchCode)
			if strings.Contains(contentFileCode, marker) {
				//println(contentFileCode, matchCode)
				unCommentCode := UnCommentCode(marker, contentFileCode)
				errUn := utils.SaveData(path, unCommentCode)
				if errUn != nil { color.Red(errUn.Error()) }
			}
			return nil
		})
}

func CommentDir(startDir string, matchCode string, excludeDirs map[string]bool)  {
	println("CommentDir ... ")
	_ = filepath.Walk(startDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil { return err }
			firstDir := utils.GetFirstDir(path)
			if _, ok := excludeDirs[firstDir]; ok { return nil }
			if info.IsDir() { return nil}
			contentFileCode := utils.ReadFile(path)
			marker := fmt.Sprintf("/* %s", matchCode)
			if strings.Contains(contentFileCode, marker) {
				//println(contentFileCode, matchCode)
				commentedCode := CommentCode(marker, contentFileCode)
				//println(commentedCode)
				errUn := utils.SaveData(path, commentedCode)
				if errUn != nil { color.Red(errUn.Error()) }
			}
			return nil
		})
}

func main()  {
	println("start agent ")

	matchCode := "ignore for building amd64-linux"

	//fmt.Println(path)
	startDir := "../../../"
	excludeDirs := map[string]bool{"vendor": true, ".git": true, "tmp": true, ".idea": true}
	//UncommentDir(startDir, matchCode, excludeDirs)
	CommentDir(startDir, matchCode, excludeDirs)
	return

	agent := Agent{
		ConfigFile: "config.agent.local.yaml",
	}
	agent.Init()


	agent.Await()

}