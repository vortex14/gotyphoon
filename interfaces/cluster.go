package interfaces

import "github.com/vortex14/gotyphoon/environment"



type ClusterProjectLabels struct {
	Git     GitLabel                 `yaml:"git,omitempty"`
	Gitlab  GitlabLabel              `yaml:"gitlab,omitempty"`
	Docker  DockerLabel              `yaml:"docker,omitempty"`
	Grafana []*GrafanaConfig `yaml:"grafana,omitempty"`
}

type ClusterProject struct {
	Name   string
	//path   string
	Config string
	Labels ClusterProjectLabels
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
