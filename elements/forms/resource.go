package forms

import (
	"context"
	"fmt"
	"github.com/vortex14/gotyphoon/utils"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type Resource struct {
	*label.MetaInfo

	LOG         interfaces.LoggerInterface
	Actions     map[string] interfaces.ActionInterface
	Resources   map[string] interfaces.ResourceInterface
	Middlewares [] interfaces.MiddlewareInterface

	// /* ignore for building amd64-linux
	parentGraph    interfaces.GraphInterface
    // */

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

func (r *Resource) SetLogger(logger interfaces.LoggerInterface) interfaces.ResourceInterface {
	r.LOG = logger
	return r
}

// /* ignore for building amd64-linux

func (r *Resource) UpdateGraphLabel() {

}

func (r *Resource) BuildEdges() interfaces.ResourceGraphInterface {
	var handlerKeys []string
	allowedMethods := make(map[string] string)
	var allowedMethodsList []string
	for _, action := range r.Actions {
		handlerKeys = append(handlerKeys, action.GetHandlerPath())
		for _, method := range action.GetMethods() {
			allowedMethods[method] = method
		}

	}

	for _, method := range allowedMethods{
		allowedMethodsList = append(allowedMethodsList, method)
	}


	r.parentGraph.BuildEdges(allowedMethodsList, handlerKeys)
	return r
}

func (r *Resource) SetGraph(graph interfaces.GraphInterface) interfaces.ResourceGraphInterface {
	r.LOG.Warning("SetGraph ------>> ", r.Name)
	r.parentGraph = graph
	return r
}

func (r *Resource) HasParentGraph() bool {
	return r.parentGraph != nil
}

func (r *Resource) GetGraph() interfaces.GraphInterface {
	return r.parentGraph
}

func (r *Resource) CreateSubGraph(options *interfaces.GraphOptions) interfaces.GraphInterface {

	return r.parentGraph.AddSubGraph(options)
}

func (r *Resource) AddGraphActionNode(action interfaces.ActionGraphInterface)  {
	if utils.IsNill(r.parentGraph) { r.LOG.Error(Errors.GraphResourceNotFound.Error()); return }
	r.LOG.Debug(
		fmt.Sprintf("adding new graph node for Action - %s, %s",
			action.GetName(),
			action.GetHandlerPath()),
	)

	action.SetGraph(r.parentGraph, true)

}

func (r *Resource) GetGraphNodes() map[string]interfaces.NodeInterface  {
	return r.parentGraph.GetNodes()
}


func (r *Resource) SetGraphNodes(nodes map[string]interfaces.NodeInterface) interfaces.ResourceGraphInterface {
	r.parentGraph.SetNodes(nodes)

	return r
}

func (r *Resource) SetGraphEdges(edges map[string]interfaces.EdgeInterface) interfaces.ResourceGraphInterface {
	r.parentGraph.SetEdges(edges)
	return r
}

func (r *Resource) GetGraphEdges() map[string]interfaces.EdgeInterface {
	return r.parentGraph.GetEdges()
}


// */