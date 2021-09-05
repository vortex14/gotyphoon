package forms

import (
	"context"
	"fmt"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"

	Errors "github.com/vortex14/gotyphoon/errors"
)


type Resource struct {
	*label.MetaInfo

	LOG         *logrus.Entry
	Actions     map[string] interfaces.ActionInterface
	Resources   map[string] interfaces.ResourceInterface
	Middlewares [] interfaces.MiddlewareInterface
}

func (r *Resource) GetActions() map[string] interfaces.ActionInterface {
	return r.Actions
}

func (r *Resource) GetResources() map[string] interfaces.ResourceInterface {
	return r.Resources
}

func (r *Resource) Get() interfaces.ResourceInterface {
	return r
}

func (r *Resource) GetCountSubResources() int {
	return len(r.Resources)
}

func (r *Resource) GetCountActions() int {
	return len(r.Actions)
}

func (r *Resource) HasResource(path string) (bool, interfaces.ResourceInterface) {
	var found bool
	var resource interfaces.ResourceInterface
	if foundResource, ok := r.Resources[path]; ok { found = true; resource = foundResource }
	return found, resource
}

func (r *Resource) HasAction(path string) (bool, interfaces.ActionInterface) {
	var found bool
	var action interfaces.ActionInterface
	if a, ok := r.Actions[path]; ok { found = true; action = a }
	return found, action
}


func (r *Resource) RunMiddlewareStack(
	ctx context.Context,
	reject func(err error),

	) {
	var failed bool

	for _, middleware := range r.Middlewares {
		if failed { break }
		logger :=  log.New(log.D{"middleware": middleware.GetName(), "resource": r.GetName()})
		middleware.Pass(ctx, logger, func(err error) {
			if middleware.IsRequired() { failed = true; reject(err) } else {
				logrus.Warning(err.Error())
			}
		}, func(context context.Context) {

		})
	}
}

func (r *Resource) AddAction(action interfaces.ActionInterface) interfaces.ResourceInterface {
	if found := r.Actions[action.GetPath()]; found != nil { color.Red("%s", Errors.ActionAlreadyExists.Error()) }
	logrus.Info(fmt.Sprintf("Registered new action <%s> for resource: < %s > ", action.GetPath(), r.GetName()))
	r.Actions[action.GetPath()] = action
	return r
}

func (r *Resource) AddResource(resource interfaces.ResourceInterface) interfaces.ResourceInterface {
	if r.Resources == nil { r.Resources = make(map[string]interfaces.ResourceInterface) }
	if found := r.Resources[resource.GetPath()]; found != nil { color.Red("%s", Errors.ResourceAlreadyExist.Error()) }
	r.Resources[resource.GetPath()] = resource
	return r
}


