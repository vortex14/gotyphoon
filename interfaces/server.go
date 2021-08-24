package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	Errors "github.com/vortex14/gotyphoon/errors"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

type ServerMetaData interface {
	GetName() string
	GetDescription() string
}

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

type MiddlewareInterface interface {
	ServerMetaData
}

type Controller func(logger *logrus.Entry, ctx *gin.Context)

type Action struct {
	Name string
	Methods []string
	Description string
	Controller Controller
	PyController Controller
}

type ActionInterface interface {
	AddMethod(name string) error
	ServerMetaData
}

type Resource struct {
	Path string
	Name string
	Description string
	Middlewares []*Middleware
	Actions map[string]*Action
	Resource map[string]*Resource
}

func (r *Resource) GetName() string {
	return r.Name
}

func (r *Resource) Get() *Resource {
	return r
}

func (r *Resource) GetDescription() string {
	return r.Description
}


func (r *Resource) AddAction(action *Action) error {
	if found := r.Actions[action.Name]; found != nil { return Errors.ActionAlreadyExists }
	r.Actions[action.Name] = action
	return nil
}


type ResourceInterface interface {
	AddAction(action *Action) error
	Get() *Resource
	ServerMetaData
}

type ServerBuilderInterface interface {
	Run(project Project) ServerInterface
}

type ServerInterface interface {
	Run() error
	Stop() error
	Restart() error
	AddResource(resource *Resource) error
	Serve(method string, path string, callback func(ctx *gin.Context))

	Init() ServerInterface
	InitDocs() ServerInterface
	InitTracer() ServerInterface
	InitLogger() ServerInterface

}

type HostServer interface {
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
	GetRunningClusters() [] Cluster
	DeployProject(project Project) error
	DeployCluster(cluster Cluster) error
}

type HostGroup interface {
	CheckNodesHealth() 	bool
	GetServers() 		[]HostServer
	GetActiveServers() 	[]HostServer
	RestartServers		([]HostServer)
	StopServers			([]HostServer)

}