package interfaces

import (
	"github.com/gin-gonic/gin"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

type MetaDataInterface interface {
	GetName() string
	GetDescription() string
}

type BaseServerLabel struct {
	Name string
	Description string
}


type BuilderInterface interface {
	Run(project Project) ServerInterface
}

type ServerInterface interface {
	Run() error
	Stop() error
	Restart() error
	AddResource(resource ResourceInterface) ServerInterface
	Serve(method string, path string, callback func(ctx *gin.Context))

	Init() ServerInterface
	InitDocs() ServerInterface
	InitTracer() ServerInterface
	InitLogger() ServerInterface

}