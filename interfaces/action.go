package interfaces

import "strings"

type Action struct {
	Name           string
	Description    string
	Path           string
	Methods        [] string  // just yet HTTP Methods
	AllowedMethods [] string
	Controller     Controller // Controller of Action
	Pipeline       PipelineGroupInterface
	PyController   Controller // Python Controller Bridge of Action
	Middlewares    [] MiddlewareInterface // Before a call to action we need to check this into middleware. May be client state isn't ready for serve
}

func (a *Action) AddMethod(name string) {
	switch name { case POST, GET, PUT, PATCH, DELETE: a.Methods = append(a.Methods, name) }
}

type ActionInterface interface {
	GetPath() string
	GetMethods() []string
	AddMethod(name string)
	GetController() Controller
	GetMiddlewareStack() []MiddlewareInterface
	GetPipeline() PipelineGroupInterface

	MetaDataInterface
}

func (a *Action) GetName() string {
	return a.Name
}

func (a *Action) GetMiddlewareStack() [] MiddlewareInterface {
	return a.Middlewares
}

func (a *Action) GetDescription() string {
	return a.Description
}

func (a *Action) GetMethods() []string {
	return a.Methods
}

func (a *Action) GetController() Controller {
	return a.Controller
}

func (a *Action) GetPipeline() PipelineGroupInterface {
	return a.Pipeline
}

func (a *Action) GetPath() string {
	return strings.ToLower(a.Path)
}