package interfaces

import (
	"context"
)

type BasePipelineInterface interface {
	Run(
		context context.Context,
		reject func(pipeline BasePipelineInterface, err error),
		next func(ctx context.Context),
	)
	RunMiddlewareStack(
		context context.Context,
		reject func(middleware MiddlewareInterface, err error),
		next func(ctx context.Context),
	)
	Cancel(
		context context.Context,
		logger LoggerInterface,
		err error,
	)
	MetaDataInterface
}

type ProcessorPipelineInterface interface {
	BasePipelineInterface
	Crawl()
	Finish()
	Switch()
}

type PipelineGroupGraph interface {
	SetGraph(graph GraphInterface)
	InitGraph(parentNode string)
	SetGraphNodes(nodes map[string]NodeInterface)

	PipelineGroupInterface
}

type PipelineGroupInterface interface {
	Run(ctx context.Context)
	GetName() string
	GetFirstPipelineName() string
	SetLogger(logger LoggerInterface)

	SetGraph(graph GraphInterface)
	InitGraph(parentNode string)
	SetGraphNodes(nodes map[string]NodeInterface)
}


type CallbackPipelineInterface interface {
	Call(ctx context.Context, data interface{})
}



type ConsumerInterface interface {

}

type LambdaInterface interface {

}

type HandlerInterface interface {
	
}

type ResponseInterface interface {
	
}