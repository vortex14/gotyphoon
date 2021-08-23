package interfaces

import (
	"github.com/vortex14/gotyphoon/builders"
	"github.com/vortex14/gotyphoon/environment"
	"github.com/vortex14/gotyphoon/migrates"
)

type Project interface {
	Run()
	Close()
	Watch()
	goPromise
	CheckProject()
	IsDebug() bool
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
	LoadServices(opts TyphoonIntegrationsOptions)
	GetService(name string) Service
}