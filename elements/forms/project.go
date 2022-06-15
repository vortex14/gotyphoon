package forms

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/folder"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/extensions/data"

	"github.com/vortex14/gotyphoon/environment"
	"github.com/vortex14/gotyphoon/integrations/mongo"
	"github.com/vortex14/gotyphoon/integrations/redis"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/services"
)

type components = struct {
	PathProject      string
	TyphoonPath      string
	ConfigFile       string
	ActiveComponents map[string]interfaces.Component
}

type Services struct {
	Mongo map[string]mongo.Service
	Redis map[string]redis.Service
}

type TestMongo struct {
}

type Project struct {
	Folder *folder.Folder
	singleton.Singleton
	awaitabler.Object

	LOG interfaces.LoggerInterface

	AutoReload        bool
	Path              string
	Name              string
	Tag               string
	LogLevel          string
	DockerImageName   string
	ConfigFile        string
	Version           string
	SelectedComponent []string
	Components        components
	repo              *git.Repository
	Watcher           fsnotify.Watcher
	Services          *services.Services
	EnvSettings       *environment.Settings
	//Archon            ghosts.ArchonInterface
	Config         *interfaces.ConfigProject
	BuilderOptions *interfaces.BuilderOptions
	Labels         *interfaces.ClusterProjectLabels
}

func (p *Project) GetDockerImageName() string {
	if len(p.DockerImageName) == 0 {
		return "typhoon-lite"
	}
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
	if p.Components.ActiveComponents == nil {
		p.InitComponents()
	}
	p.Components.ActiveComponents["fetcher"].InitConsumers(p)
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
	return nil
}

func (p *Project) Watch() {

}

func (p *Project) Close() {
}

func (p *Project) Down() {
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

}

func (p *Project) InitComponents() {

}

func (p *Project) StartComponents(promise bool) {

}

func (p *Project) GetVersion() string {
	return p.Version
}

func (p *Project) CreateSymbolicLink() error {
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
}

func (p *Project) LoadConfig() (configProject *interfaces.ConfigProject) {
	if p.Config != nil {
		return p.Config
	}
	configPath := fmt.Sprintf("%s/%s", p.GetProjectPath(), p.ConfigFile)
	color.Yellow("Load config from file: %s", configPath)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		color.Red("Config %s does not exists in project :%s", p.ConfigFile, configPath)
		panic("Config does not exists in project")
	}

	var loadedConfig interfaces.ConfigProject
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("config.yaml err   #%v ", err)
		panic("config err")
	} else {
		err = yaml.Unmarshal(yamlFile, &loadedConfig)
		if err != nil {
			//log.Fatalf("Unmarshal: %v", err)
			color.Red("Config load error: %s", err)
			panic("Config Unmarshal error")
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
}

func (p *Project) RunArchon(promise bool) {
	//if p.Archon == nil {
	//	color.Red("Archon doesn't exist in your project")
	//	os.Exit(1)
	//}
	//p.LoadConfig()
	//p.Archon.RunDemons(p)
	//
	//p.Archon.RunProjectServers(p)
	//
	//if promise {
	//	p.Archon.AddPromise()
	//}

}
