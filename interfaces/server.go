package interfaces

const (
	GET     = "GET"
	POST    = "POST"
	OPTIONS = "OPTIONS"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
)

type Response map[string]interface{}

type ManageServerInterface interface {
	Run() error
	Stop() error
	Restart() error
	GetServerEngine() interface{}
	SetRouterGroup(resource ResourceInterface, group interface{})
}

// /* ignore for building amd64-linux
//
//type ServerGraphInterface interface {
//	InitGraph() ServerInterface
//	GetGraph() GraphInterface
//}
//
// */

type ServerExtensionInterface interface {
	RunServer(port int) error
	InitResourcesMap()
}

type BuilderInterface interface {
	Run(project Project) ServerInterface
}

type ServerOptionsInterface interface {
	Init() ServerInterface
	InitDocs() ServerInterface
	InitTracer() ServerInterface
	InitLogger() ServerInterface
}

type ServerInterface interface {
	AddResource(resource ResourceInterface) ServerInterface

	ServerExtensionInterface
	ServerOptionsInterface
	ManageServerInterface

	MetaDataInterface
}