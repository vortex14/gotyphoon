package interfaces

import "context"

type ActionGraphInterface interface {

	///* ignore for building amd64-linux

	SetGraph(graph GraphInterface, buildMethods bool)
	SetGraphNodes(nodes map[string] NodeInterface)
	UpdateGraphLabel(method string, path string)

	//*/

	InitPipelineGraph()

	ActionInterface
}

type ActionInterface interface {
	GetPath() string
	GetMethods() []string
	SetHandlerPath(path string)
	SetLogger(logger LoggerInterface)
	GetHandlerPath() string
	AddMethod(name string)

	OnRequest(method string, path string)

	GetController() Controller
	IsPipeline() bool
	GetMiddlewareStack() []MiddlewareInterface
	GetPipeline() PipelineGroupInterface

	Run(ctx context.Context, logger LoggerInterface)
	MetaDataInterface


}