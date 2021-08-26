package server

import "github.com/vortex14/gotyphoon/interfaces"

type Host interface {
	RunCommand()
	GetSSHClient()
	Restart() error
	StopAllClusters()
	StopAllProjects()
	UpdateTyphoonNode()
	StopSystemdService()
	RunAnsiblePlaybook()
	PrepareTyphoonNode()
	CheckFreeDiskSpace()
	CreateSystemdService()
	CheckNodeHealth() bool
	CreateSSHAccessUserRecord()
	GetRunningClusters() []interfaces.Cluster
	DeployProject(project interfaces.Project) error
	DeployCluster(cluster interfaces.Cluster) error
}

type HostGroup interface {
	CheckNodesHealth() 	bool
	GetServers() 		[]Host
	GetActiveServers() 	[]Host
	RestartServers		([]Host)
	StopServers			([]Host)

}
