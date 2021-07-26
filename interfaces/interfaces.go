package interfaces

import (
	"bufio"
	"github.com/vortex14/gotyphoon/builders"
	"github.com/vortex14/gotyphoon/environment"
	"github.com/vortex14/gotyphoon/migrates"
	"os"
)


type GitlabServer interface {
	GetAllProjectsList() []*GitlabProject
	SyncGitlabProjects()
	Deploy()
	HistoryPipelines()
}


type ConfigInterface interface {
	GetComponentPort(name string) int
	SetConfigName(name string)
	GetConfigName() string
	SetConfigPath(path string)
	GetConfigPath() string
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


type QueueSettings interface {
	SetGroupName(name string)
	GetGroupName() string
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

type Group interface {
	CheckNodesHealth() bool
	GetServers() []*Server
	GetActiveServers() []*Server
	RestartServers([]*Server)
	StopServers([]*Server)

}


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
	GetHost() string
	GetPort() int
}

type AdapterService interface {
	Ping() bool
	Init()
}

type Database interface {
	Import(Database string, collection string, inputFile string) (error, uint64)
	Export(Database string, collection string, outFile string) (*bufio.Writer, *os.File, int64, error)
}


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
	LoadServices()
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
	LoadConfig() *ConfigProject
	GetSelectedComponent() []string
	GetComponentPort(name string) int
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






