package interfaces

import (
	"github.com/vortex14/gotyphoon/builders"
	"github.com/vortex14/gotyphoon/config"
	"github.com/vortex14/gotyphoon/environment"
	"github.com/vortex14/gotyphoon/migrates"
	"github.com/xanzy/go-gitlab"
)

type GoTemplate struct {
	Source string
	ExportPath string
	Data interface{}
}


type ClusterLabel struct {
	Kind string
	Version string
}

type ClusterGitlab struct {
	Endpoint string `yaml:"endpoint,omitempty"`
	Variables []*gitlab.PipelineVariable `yaml:"variables,omitempty"`
}

type ClusterGrafana struct {
	Endpoint string `yaml:"endpoint,omitempty"`
	FolderId string `yaml:"folder_id,omitempty"`
}

type ClusterDocker struct {
	Image string
}

type ClusterMeta struct {
	Gitlab  ClusterGitlab  `yaml:"gitlab,omitempty"`
	Grafana ClusterGrafana `yaml:"grafana,omitempty"`
	Docker  ClusterDocker  `yaml:"docker,omitempty"`
}

type GitlabLabel struct {
	Id int `yaml:"id,omitempty"`
}

type GitlabClusterLabel struct {
	Url int `yaml:"url,omitempty"`
}

type GitLabel struct {
	Url string `yaml:"url,omitempty"`
	Remote string `yaml:"remote,omitempty"`
	Branch string `yaml:"branch,omitempty"`
}

type DockerLabel struct {

}

type FileObject struct {
	Type string
	Path string
	Name string
	Data string
	FileSystem
}

type ClusterProjectLabels struct {
	Git     GitLabel                 `yaml:"git,omitempty"`
	Gitlab  GitlabLabel              `yaml:"gitlab,omitempty"`
	Docker  DockerLabel              `yaml:"docker,omitempty"`
	Grafana []*config.GrafanaConfig `yaml:"grafana,omitempty"`
}

type ClusterProject struct {
	Name   string
	path   string
	Config string
	Labels ClusterProjectLabels
}

type GitlabProject struct {
	Name string `yaml:"name,omitempty"`
	Git string `yaml:"git,omitempty"`
	Id int	`yaml:"id,omitempty"`
}

type GitlabServer interface {
	GetAllProjectsList() []*GitlabProject
	SyncGitlabProjects()
	Deploy()
	HistoryPipelines()
}

type Cluster interface {
	Add()
	Show()
	Create()
	Deploy()
	SaveConfig()
	GetName() string
	GetConfigName() string
	GetClusterConfigPath() string
	GetProjects() [] *ClusterProject
	GetMeta() *ClusterMeta
	GetEnvSettings() *environment.Settings
}


type Server interface {
	CheckNodeHealth() bool
	Restart() error
	GetSSHClient()
	DeployCluster(cluster *Cluster) error
	DeployProject(project *Project) error
	RunCommand()
	GetRunningClusters() []*Cluster
	CreateSystemdService()
	StopSystemdService()
	RunAnsiblePlaybook()
	CreateSSHAccessUserRecord()
	PrepareTyphoonNode()
	UpdateTyphoonNode()
	StopAllProjects()
	StopAllClusters()
	CheckFreeDiskSpace()
}


type K8sCluster interface {
	PortForward()
}

type CloudManagement interface {
	Deploy()
}

type Producer interface {
	Pub()
}

type Consumer interface {
	Read()
}

type Pipeline interface {
	Run()
	Crawl()
	Finish()
	Switch()
	Retry()
	Await()
}



type Producers map[string] Producer
type Consumers map[string] [] Consumer
type Pipelines map[string] Pipeline

type Group interface {
	CheckNodesHealth() bool
	GetServers() []*Server
	GetActiveServers() []*Server
	RestartServers([]*Server)
	StopServers([]*Server)

}

type ReplaceLabel struct {
	Label string
	Value string
}

type ReplaceLabels []*ReplaceLabel


type MapFileObjects map[string]*FileObject
type BuilderOptions builders.BuildOptions

type FileSystem interface {
	GetDataFromDirectory(path string) MapFileObjects
	IsExistDir (path string) bool
}

type TestData interface {
	GetFields()
}

type Environment interface {
	Load()
	Set()
	Get()
	GetSettings() (error, *environment.Settings)
}

type goPromise interface {
	AddPromise()
	PromiseDone()
	WaitPromises()
}


type Service interface {
	TestConnect() bool
}

type Services struct {}

type GrafanaInterface interface {
	ImportGrafanaConfig()
	RemoveGrafanaDashboard()
	CreateBaseGrafanaConfig()
	CreateGrafanaMonitoringTemplates()
}


type DockerInterface interface {
	BuildImage()
	ListContainers()
	ProjectBuild()
	RunComponent(component string) error
}

type HelmInterface interface {
	BuildHelmMinikubeResources()
	RemoveHelmMinikubeManifests()
}

type GitlabInterface interface {
	BuildCIResources()
}

type Project interface {
	Run()
	Close()
	Watch()
	CheckProject()
	GetTag() string
	GetName() string
	GetVersion() string
	GetLogLevel() string
	GetConfigFile() string
	GetProjectPath() string
	builders.ProjectBuilder
	migrates.ProjectMigrate
	GetComponents() []string
	CreateSymbolicLink() error
	GetDockerImageName() string
	GetComponentPort(name string) int
	LoadConfig() *config.Project
	GetLabels() *ClusterProjectLabels
	GetBuilderOptions() *BuilderOptions
	GetEnvSettings() *environment.Settings
	goPromise
}

type Utils interface {
	GoRunTemplate(goTemplate *GoTemplate) bool
	ParseLog(object *FileObject) error
	GetGoTemplate(object *FileObject) (error, string)
}

type Component interface {
	CheckDirectory(required []string, pathComponent string) bool
	GetName() string
	//CheckComponent(component string) bool
	Start(project Project)
	Close(project Project)
	Stop(project Project)
	//Restart(project *Project)
	goPromise
}






