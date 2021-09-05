package interfaces

import "context"

type ActionInterface interface {
	GetPath() string
	GetMethods() []string
	AddMethod(name string)

	GetController() Controller
	GetMiddlewareStack() []MiddlewareInterface
	GetPipeline() PipelineGroupInterface

	Run(ctx context.Context, logger LoggerInterface)
	MetaDataInterface
}