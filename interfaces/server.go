package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

type BaseServerLabel struct {
	Name string
	Description string
}

type Middleware struct {
	Name string
	Description string
	ctx *gin.Context
	Callback func(ctx *gin.Context)
	PyCallback func(ctx *gin.Context)

}

type Controller func(logger *logrus.Entry, ctx *gin.Context)

type Action struct {
	Name string
	Methods []string
	Description string
	Controller Controller
	PyController Controller

}

type Resource struct {
	Path string
	Name string
	Description string
	Middlewares []*Middleware
	Actions map[string]*Action
	Resource map[string]*Resource
}

type ServerInterface interface {
	Run() error
	Stop() error
	Restart() error
	Serve(method string, path string, callback func(ctx *gin.Context))
	AddResource(resource *Resource) error
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

type Group interface {
	CheckNodesHealth() 	bool
	GetServers() 		[]Server
	GetActiveServers() 	[]Server
	RestartServers		([]Server)
	StopServers			([]Server)

}