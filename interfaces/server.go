package interfaces

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

type Response map[string]interface{}

type ManageServerInterface interface {
	Run() error
	Stop() error
	Restart() error
}

type ServerGraphInterface interface {
	InitGraph() ServerInterface
	//GetGraph() GraphInterface
}

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