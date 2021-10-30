package interfaces

import (
	"context"
)


type ResourceGraphInterface interface {

	// /* ignore for building amd64-linux

	SetGraph(graph GraphInterface) ResourceGraphInterface
	GetGraph() GraphInterface

	CreateSubGraph(options *GraphOptions) GraphInterface
	AddGraphActionNode(action ActionGraphInterface)

	GetGraphNodes() map[string] NodeInterface
	SetGraphNodes(nodes map[string] NodeInterface) ResourceGraphInterface

	BuildEdges() ResourceGraphInterface
	SetGraphEdges(edges map[string]EdgeInterface) ResourceGraphInterface
	GetGraphEdges()map[string] EdgeInterface

	HasParentGraph() bool

	// */

	ResourceInterface
}

type ResourceAuthInterface interface {
	Allow(server ServerInterface, resource ResourceInterface) interface{}
	SetServerEngine(server ServerInterface)
	SetLogger(logger LoggerInterface)
}

type ResourceInterface interface {
	GetPath() string
	GetCountActions()int
	Get() ResourceInterface
	SetRouterGroup(group interface{})
	GetRouterGroup() interface{}
	GetCountSubResources() int
	SetDescription(description string)
	GetActions() map[string] ActionInterface
	GetResources() map[string] ResourceInterface
	IsAuth() bool
	InitAuth(server ServerInterface)

	SetLogger(logger LoggerInterface) ResourceInterface

	HasAction(path string) (bool, ActionInterface)
	HasResource(path string) (bool, ResourceInterface)

	SetName(name string)
	RunMiddlewareStack(context context.Context, reject func(err error))

	AddAction(action ActionInterface) ResourceInterface
	AddResource(resource ResourceInterface) ResourceInterface


	MetaDataInterface
}