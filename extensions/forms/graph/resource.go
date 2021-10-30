package graph

// /* ignore for building amd64-linux

import (
	"fmt"
	"github.com/vortex14/gotyphoon/elements/forms"

	"github.com/vortex14/gotyphoon/utils"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Resource struct {
	*forms.Resource

	parentGraph    interfaces.GraphInterface
}

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