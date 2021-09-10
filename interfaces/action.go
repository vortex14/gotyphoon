package interfaces

import "context"

type ActionInterface interface {
	GetPath() string
	GetMethods() []string
	SetHandlerPath(path string)
	GetHandlerPath() string
	AddMethod(name string)

	GetController() Controller
	GetMiddlewareStack() []MiddlewareInterface
	GetPipeline() PipelineGroupInterface
	UpdateGraphLabel(method string, path string)
	Run(ctx context.Context, logger LoggerInterface)
	MetaDataInterface

	SetGraph(graph GraphInterface)
}