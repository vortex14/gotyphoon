package interfaces

import (
	"context"
)

type ResourceInterface interface {
	GetPath() string
	GetCountActions()int
	Get() ResourceInterface
	GetCountSubResources() int
	SetDescription(description string)
	GetActions() map[string] ActionInterface
	GetResources() map[string] ResourceInterface

	SetLogger(logger LoggerInterface) ResourceInterface

	HasAction(path string) (bool, ActionInterface)
	HasResource(path string) (bool, ResourceInterface)

	SetName(name string)
	RunMiddlewareStack(context context.Context, reject func(err error))

	AddAction(action ActionInterface) ResourceInterface
	AddResource(resource ResourceInterface) ResourceInterface

	SetGraph(graph GraphInterface) ResourceInterface
	AddGraphActionNode(action ActionInterface)

	MetaDataInterface
}