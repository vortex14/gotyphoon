package server

import (
	"github.com/fatih/color"
	Errors "github.com/vortex14/gotyphoon/errors"
)

type ResourceInterface interface {
	AddAction(action *Action) ResourceInterface
	AddResource(resource *Resource) ResourceInterface
	Get() *Resource
	MetaDataInterface
}

type Resource struct {
	Path string
	Name string
	Description string
	Middlewares []*Middleware
	Actions map[string]*Action
	Resource map[string]*Resource
}

func (r *Resource) GetName() string {
	return r.Name
}

func (r *Resource) Get() *Resource {
	return r
}

func (r *Resource) GetDescription() string {
	return r.Description
}


func (r *Resource) AddAction(action *Action) ResourceInterface {
	if found := r.Actions[action.Name]; found != nil { color.Red("%s", Errors.ActionAlreadyExists.Error()) }
	r.Actions[action.Name] = action
	return r
}

func (r *Resource) AddResource(resource *Resource) ResourceInterface  {
	if found := r.Resource[resource.Path]; found != nil { color.Red("%s", Errors.ResourceAlreadyExist.Error()) }
	r.Resource[resource.Path] = resource
	return r
}

