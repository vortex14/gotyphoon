package server

import (
	"github.com/gin-gonic/gin"

	"github.com/vortex14/gotyphoon/interfaces"
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
	Run(project interfaces.Project) Interface
}

type Interface interface {
	Run() error
	Stop() error
	Restart() error
	AddResource(resource *Resource) Interface
	Serve(method string, path string, callback func(ctx *gin.Context))

	Init() Interface
	InitDocs() Interface
	InitTracer() Interface
	InitLogger() Interface

}