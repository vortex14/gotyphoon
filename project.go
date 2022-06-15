package typhoon

import (
	"bufio"
	"context"
	"fmt"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"crypto/md5"
	"encoding/hex"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/folder"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/extensions/data"

	Errors "github.com/vortex14/gotyphoon/errors"
	tyLog "github.com/vortex14/gotyphoon/log"

	"github.com/vortex14/gotyphoon/environment"
	"github.com/vortex14/gotyphoon/integrations/mongo"
	"github.com/vortex14/gotyphoon/integrations/redis"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/interfaces/ghosts"
	"github.com/vortex14/gotyphoon/services"
)

type components = struct {
	PathProject      string
	TyphoonPath      string
	ConfigFile       string
	ActiveComponents map[string]*Component
}

type Task struct {
	ticker *time.Ticker
	closed chan struct{}
	wg     sync.WaitGroup
}

type Services struct {
	Mongo map[string]mongo.Service
	Redis map[string]redis.Service
}

type TestMongo struct {
}

type Project struct {
	folder *folder.Folder
	singleton.Singleton
	awaitabler.Object

	LOG interfaces.LoggerInterface

	AutoReload        bool
	task              *Task
	Path              string
	Name              string
	Tag               string
	LogLevel          string
	DockerImageName   string
	ConfigFile        string
	Version           string
	SelectedComponent []string
	components        components
	repo              *git.Repository
	Watcher           fsnotify.Watcher
	Services          *services.Services
	EnvSettings       *environment.Settings
	Archon            ghosts.ArchonInterface
	Config            *interfaces.ConfigProject
	BuilderOptions    *interfaces.BuilderOptions
	Labels            *interfaces.ClusterProjectLabels
}

func (p *Project) GetDockerImageName() string {
	return p.DockerImageName
}

func (p *Project) GetLabels() *interfaces.ClusterProjectLabels {
	return p.Labels
}

func (p *Project) IsDebug() bool {

	return p.Config.Debug
}
func (p *Project) RunFetcherQueues() {
	p.LoadConfig()
	if p.components.ActiveComponents == nil {
		p.initComponents()
	}
	p.components.ActiveComponents["fetcher"].InitConsumers(p)
}

func (p *Project) GetService(name string) interfaces.Service {
	switch name {
	case interfaces.NSQ:
		return p.Services.Collections.Nsq
		//case interfaces.MONGO:
		//	return p.Services.Collections.Mongo

	}
	return nil
}

func (p *Project) GetRepo() (error, *git.Repository) {
	if p.repo == nil {
		repo, err := git.PlainOpen(p.GetProjectPath())
		if err != nil {
			return err, nil
		}
		p.repo = repo
	}
	return nil, p.repo
}

func (p *Project) GetConfigs() []string {
	var configs []string
	re := regexp.MustCompile(`config.*.yaml`)
	files, errR := ioutil.ReadDir(p.GetProjectPath())
	if errR != nil {
		color.Red("%s,", errR.Error())
		return nil
	}
	for _, f := range files {
		found := string(re.Find([]byte(f.Name())))
		if len(found) == 0 {
			continue
		}
		configs = append(configs, found)
	}
	return configs
}

func (p *Project) GetBranch() (error, string) {
	var branchName string
	var err error
	errRepo, repo := p.GetRepo()
	repoData, errHead := repo.Head()
	if errHead != nil {
		err = errHead
	}
	if errRepo != nil {
		err = errRepo
	}

	if repoData != nil {
		branchName = repoData.Name().Short()
	}

	return err, branchName
}

func (p *Project) GetRemotes() ([]*git.Remote, error) {
	var err error
	errRepo, repo := p.GetRepo()
	remotes, errRemote := repo.Remotes()
	if errRepo != nil {
		err = errRepo
	}
	if errRemote != nil {
		err = errRemote
	}

	return remotes, err
}

func watchDirTeet(path string, fi os.FileInfo, err error) error {
	print(path)

	return nil
	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}

func (p *Project) GetComponentPort(name string) int {
	return p.Config.GetComponentPort(name)
}

func (p *Project) WatchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return p.Watcher.Add(path)
	}

	return nil
}

func (p *Project) RunTestServices() {

	typhoonServices := services.Services{Project: p}
	typhoonServices.RunTestServices()
}

func (p *Project) ImportExceptions(component string, sourceFileName string) error {
	currentPath, _ := os.Getwd()
	importPath := fmt.Sprintf("%s/%s", currentPath, sourceFileName)
	f, err := os.OpenFile(importPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		color.Red("open file error: %v", err)
		return err
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			return err
		}

		//BsonTools.m
		//var doc BsonTools.BSON
		//bson.UnmarshalExtJSON()

		//doc = suurceBson
		//
		//bson.MarshalExtJSONAppend()

		//json, err := doc.JSON()
		//if err != nil {
		//	color.Red("--- > %+v", err)
		//	return err
		//}
		var doc bson.D

		if err := bson.UnmarshalExtJSON([]byte(line), false, &doc); err != nil {
			return err
		}

		//d, _ := bson.Ma(doc)

		//u := utils.Utils{}
		//j

		//o := u.DumpPrettyJson(suurceBson)

		//println(BsonTools)

		//color.Green("%+v", d)  // GET the line string

	}

	return nil
}

func (p *Project) ImportResponseData(url string, sourceFile string) {
	p.LoadConfig()

	currentPath, _ := os.Getwd()
	importPath := fmt.Sprintf("%s/%s", currentPath, sourceFile)
	dat, err := ioutil.ReadFile(importPath)
	if err != nil {
		color.Red("%s", err)
		os.Exit(1)
	}
	color.Green("url: %s", url)
	taskid := md5.Sum([]byte(url))
	p.LoadServices(interfaces.TyphoonIntegrationsOptions{
		Redis: interfaces.BaseServiceOptions{
			Active: true,
		},
	},
	)
	redisPath := fmt.Sprintf("%s:%s", p.GetName(), hex.EncodeToString(taskid[:]))
	err = p.Services.Collections.Redis["main"].Set(redisPath, string(dat))
	if err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}
	color.Green(redisPath)
}

func (p *Project) TestFunc() {
	data.TestFunc()
}

func (p *Project) CreateProject() {

}

func (p *Project) BuildCIResources() {

}

func (p *Project) GetEnvSettings() *environment.Settings {
	return p.EnvSettings
}

func (p *Project) AddPromise() {
	p.Add()
}
func (p *Project) PromiseDone() {
	p.Done()
}
func (p *Project) WaitPromises() {
	p.Await()
}
func (p *Project) Run() interfaces.Project {
	p.folder = &folder.Folder{Path: p.GetProjectPath()}

	p.CheckProject()

	p.LOG = tyLog.New(tyLog.D{"project": p.Name})

	//println(">>>>>>>> >>>>>>>", p.folder.Path)

	p.LOG.Info(p.folder.Path)

	if !p.folder.IsExist("typhoon") {
		p.LOG.Info("init typhoon symlink path")
		_ = p.CreateSymbolicLink()
	}

	//
	//
	//color.Magenta("start components")
	p.LOG.Info("start components ...")
	//p.AddPromise()
	p.StartComponents(true)

	return p

}

func (p *Project) Watch() {
	color.Green("watch for project ..")
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	if err := filepath.Walk("project", watchDirTeet); err != nil {
		fmt.Println("ERROR", err)
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:

				if strings.Contains(event.Name, ".pyc") {
					continue
				}

				if strings.Contains(event.String(), "CHMOD") {
					continue
				}

				if strings.Contains(event.Name, ".py~") {
					continue
				}

				if strings.Contains(event.Name, "__pycache__") {
					continue
				}

				componentChanged := ""

				for _, component := range p.SelectedComponent {
					if strings.Contains(event.Name, strings.ToLower(component)) {
						color.Yellow("reloading %s ... !", component)
						componentChanged = component
						break
					}

				}

				if _, ok := p.components.ActiveComponents[componentChanged]; ok {

					color.Yellow("Reload %s ...", componentChanged)
					color.Yellow("event %+v", event)
					component := p.components.ActiveComponents[componentChanged]

					//p.AddPromise()
					go component.Restart(p)

					// "example" is not in the map
				} else {
					color.Yellow("%s isn't running", componentChanged)
				}

				//

				//p.AddPromise()
				//go component.Restart(p)

				//go component.Restart(p)

				//initComponent(wg, tcomponents, componentChanged, configFile)

				// watch for errors
			case err := <-watcher.Errors:
				color.Red("ERROR---->", err)
			}
		}
	}()

	<-done
}

func (p *Project) Close() {
	color.Yellow("close project ...")
	for _, component := range p.components.ActiveComponents {
		if component.Active {
			p.Add()
			go component.Close(p)
		}

	}
	color.Yellow("await done...")
	p.Await()
	color.Yellow("await done !")

}

func (p *Project) Down() {
	p.LoadConfig()
	commandDropProject := fmt.Sprintf("kill -9 $(ps aux | grep \"%s\" | awk '{print $2}')", p.GetName())
	color.Red("Running: %s: ", commandDropProject)
	ctxP, cancelP := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelP()

	if err := exec.CommandContext(ctxP, "bash", "-c", commandDropProject).Run(); err != nil {
		color.Yellow("Project components killed!")
		// This will fail after 100 milliseconds. The 5 second sleep
		// will be interrupted.
	}

	commandDropTyphoon := fmt.Sprintf("kill -9 $(ps aux | grep \"%s\" | awk '{print $2}')", "typhoon")
	ctxT, cancelT := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelT()

	if err := exec.CommandContext(ctxT, "bash", "-c", commandDropTyphoon).Run(); err != nil {
		color.Red("%s", err.Error())
	}
}

func (p *Project) GetBuilderOptions() *interfaces.BuilderOptions {
	return p.BuilderOptions
}

func (p *Project) GetTag() string {
	return p.Tag
}
func (p *Project) Migrate() {

}

func (p *Project) Build() {
	color.Yellow("builder run... options %+v", p.BuilderOptions)
}

func (p *Project) GetSelectedComponent() []string {
	return p.SelectedComponent
}

func (p *Project) RunQueues() {
	if len(p.SelectedComponent) == 0 {
		color.Red("No set components for project")
		return
	}
	p.Services.RunNSQ()
}

func (p *Project) initComponents() {
	p.components.ActiveComponents = make(map[string]*Component)
	p.Name = p.Config.ProjectName
	for _, componentName := range p.SelectedComponent {

		var componentFileName string

		switch componentName {
		case interfaces.FETCHER:
			componentFileName = interfaces.TYPHOON2PYTHON2FETCHER
		case interfaces.PROCESSOR:
			componentFileName = interfaces.TYPHOON2PYTHON2PROCESSOR
		case interfaces.DONOR:
			componentFileName = interfaces.TYPHOON2PYTHON2DONOR
		case interfaces.TRANSPORTER:
			componentFileName = interfaces.TYPHOON2PYTHON2TRANSPORTER
		case interfaces.SCHEDULER:
			componentFileName = interfaces.TYPHOON2PYTHON2SCHEDULER

		}

		p.components.ActiveComponents[componentName] = &Component{
			Component: forms.Component{
				MetaInfo: label.MetaInfo{Name: componentName},
				Language: interfaces.PYTHON,
				FileExt:  fmt.Sprintf("%s.py", componentFileName),
			},
		}
	}
}

func (p *Project) StartComponents(promise bool) {

	fmt.Printf(tyLog.DOG)

	if p.components.ActiveComponents == nil {
		p.initComponents()
	}

	//if promise {
	//	defer p.PromiseDone()
	//}

	for _, componentName := range p.SelectedComponent {
		p.Add()
		go p.components.ActiveComponents[componentName].Start(p)
	}

	p.LOG.Debug("Components started")
}

func (p *Project) GetVersion() string {
	return p.Version
}

func (p *Project) CreateSymbolicLink() error {
	env := &environment.Environment{}
	_, settings := env.GetSettings()

	linkTyphoonPath := fmt.Sprintf("%s/pytyphoon/typhoon", settings.Path)
	color.Yellow("TYPHOON_PATH=%s", settings.Path)
	directLink := filepath.Join(p.GetProjectPath(), "typhoon")
	color.Yellow(directLink)
	err := os.Symlink(linkTyphoonPath, directLink)

	if err != nil {
		fmt.Printf("err %s", err)
	}

	return nil
}

func (p *Project) GetName() string {
	projectName := p.Name
	if len(projectName) == 0 {
		projectName = p.Config.ProjectName
	}
	return projectName
}

func (p *Project) GetComponents() []string {
	return p.SelectedComponent
}

func (p *Project) GetConfigFile() string {
	return p.ConfigFile
}

func (p *Project) GetConfigPath() string {
	return filepath.Join(p.GetProjectPath(), p.GetConfigFile())
}

func (p *Project) GetProjectPath() string {
	var pathProject string
	if len(p.Path) > 0 {
		pathProject = p.Path
	} else {
		ProjectCurrent, err := os.Getwd()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		} else {
			pathProject = ProjectCurrent
		}
	}

	return pathProject
}
func (p *Project) GetLogLevel() string {
	if len(p.LogLevel) == 0 {
		config := p.LoadConfig()
		if config.Debug {
			p.LogLevel = interfaces.DEBUG
		} else {
			p.LogLevel = interfaces.INFO
		}
	}

	return p.LogLevel
}

func (p *Project) LoadServices(opts interfaces.TyphoonIntegrationsOptions) {
	status := false
	projectServices := services.Services{
		Project: p,
		Options: opts,
	}

	if opts.NSQ.Active {
		projectServices.RunNSQ()
		status = true
	}

	if opts.Mongo.Active {
		projectServices.LoadMongoServices()
		status = true
	}

	if opts.Redis.Active {
		projectServices.LoadRedisServices()
		status = true
	}
	p.Services = &projectServices

	color.Yellow("LoadServices: %t", status)
}

func (p *Project) LoadConfig() (configProject *interfaces.ConfigProject) {
	if p.Config != nil {
		return p.Config
	}
	configPath := fmt.Sprintf("%s/%s", p.GetProjectPath(), p.ConfigFile)
	color.Yellow("Load config from file: %s", configPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		color.Red("Config %s does not exists in project :%s", p.ConfigFile, configPath)
		os.Exit(1)
	}

	var loadedConfig interfaces.ConfigProject
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("config.yaml err   #%v ", err)
		os.Exit(1)
	} else {
		err = yaml.Unmarshal(yamlFile, &loadedConfig)
		if err != nil {
			//log.Fatalf("Unmarshal: %v", err)
			color.Red("Config load error: %s", err)
			os.Exit(1)
		}

	}

	loadedConfig.SetConfigName(p.ConfigFile)

	loadedConfig.SetConfigPath(configPath)

	//color.Yellow("Set Config details ... %s, %s", configLoad.GetConfigName(), configLoad.GetConfigPath())

	p.Config = &loadedConfig

	env := &environment.Environment{}
	_, settings := env.GetSettings()

	p.EnvSettings = settings

	return &loadedConfig
}

func (p *Project) CheckProject() {

	var status = true

	for _, componentName := range p.SelectedComponent {
		component := &Component{
			Component: forms.Component{
				MetaInfo:    label.MetaInfo{Name: componentName},
				ProjectPath: p.folder.Path,
			},
		}
		color.Yellow("checking: %s...", componentName)

		if !component.CheckComponent() {
			status = false
			color.Yellow("%s is: false", componentName)
		}

	}
	p.LoadConfig()

	if !status {
		color.Red("%s : %s", Errors.ProjectNotFound.Error(), p.Path)
		os.Exit(1)
	}

	env := &environment.Environment{}
	_, settings := env.GetSettings()

	if len(settings.Path) == 0 || len(settings.Projects) == 0 {
		color.Red("%s in %s", Errors.ProjectInvalidEnv.Error(), env.ProfilePath)
		os.Exit(1)
	}

	p.EnvSettings = settings

}

func (p *Project) RunArchon(promise bool) {
	if p.Archon == nil {
		color.Red("Archon doesn't exist in your project")
		os.Exit(1)
	}
	p.LoadConfig()
	p.Archon.RunDemons(p)

	p.Archon.RunProjectServers(p)

	if promise {
		p.Archon.AddPromise()
	}

}
