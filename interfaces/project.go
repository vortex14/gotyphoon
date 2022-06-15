package interfaces

import (
	"github.com/vortex14/gotyphoon/environment"
	//"github.com/vortex14/gotyphoon/extensions/project/python3"
	"github.com/vortex14/gotyphoon/extensions/project/python3/builders"
)

type Project interface {
	Run() Project
	Close()
	Watch()
	goPromise
	CheckProject()
	IsDebug() bool
	GetTag() string
	GetName() string
	GetVersion() string
	GetLogLevel() string
	GetConfigPath() string
	GetConfigFile() string
	GetProjectPath() string
	builders.ProjectBuilder
	//python3.ProjectMigrate
	GetComponents() []string
	CreateSymbolicLink() error
	GetDockerImageName() string
	LoadConfig() *ConfigProject
	GetSelectedComponent() []string
	GetComponentPort(name string) int
	GetLabels() *ClusterProjectLabels
	GetBuilderOptions() *BuilderOptions
	GetEnvSettings() *environment.Settings
	GetService(name string) Service
	LoadServices(opts TyphoonIntegrationsOptions)
}
