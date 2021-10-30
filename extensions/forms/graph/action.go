package graph

// /* ignore for building amd64-linux

import (
	"context"
	"fmt"
	"github.com/vortex14/gotyphoon/elements/forms"

	graphvizExt "github.com/vortex14/gotyphoon/extensions/models/graphviz"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Action struct {
	*forms.Action
	Graph          interfaces.GraphInterface
}

func (a *Action) InitPipelineGraph(logger interfaces.LoggerInterface)  {
	a.Pipeline.SetGraph(a.Graph)
	a.Pipeline.InitGraph(a.Path)
	a.Pipeline.SetGraphNodes(a.Graph.GetNodes())
}

func (a *Action) UpdateGraphLabel(method string, path string)  {
	a.Input ++

		labelAction := fmt.Sprintf(`
	
	 R: %d
	
	`, a.Input)

	a.Graph.UpdateEdge(&interfaces.EdgeOptions{
		NodeA: method,
		NodeB: path,
		LabelH: labelAction,
		Color: graphvizExt.COLORNavy,

	})
}

func (a *Action) AddMethodNodes() {
	for _, method := range a.GetMethods() {
		a.Graph.AddNode(&interfaces.NodeOptions{
			Name: method,
			Label: graphvizExt.FormatBottomSpace(method),
			Shape: graphvizExt.SHAPEAction,
			EdgeOptions: &interfaces.EdgeOptions{
				NodeB:  a.GetHandlerPath(),
				ArrowS: 0.5,
			},
		})
	}
}

func (a *Action) SetGraph(parent interfaces.GraphInterface, buildMethods bool)  {
	a.Graph = parent

	opts := &interfaces.NodeOptions{
		Name: a.GetHandlerPath(),
		Label: graphvizExt.FormatSpace(a.GetHandlerPath()),
		Shape: graphvizExt.SHAPETab,
		Style: graphvizExt.StyleFilled,
		BackgroundColor: graphvizExt.COLORGray,
		EdgeOptions: &interfaces.EdgeOptions{},

	}

	a.Graph.AddNode(opts)

	if buildMethods {
		a.AddMethodNodes()
	}

	if a.IsPipeline() {

		println(fmt.Sprintf("+%v ", a.Graph.GetNodes()), "======== !!!!!!!!!! ***********")
		println(fmt.Sprintf("!!!!!!!!!!!!!!!!!!!!!!!!!! =========== = = == == %+v", parent.GetNodes()))
		opts.EdgeOptions.NodeB = a.Pipeline.GetFirstPipelineName()
		println(opts, "++++++++++++++++++++++++", a.Pipeline.GetName(), a.IsPipeline())
	}
}

func (a *Action) SetGraphNodes(nodes map[string]interfaces.NodeInterface)  {
	println(fmt.Sprintf("ACTION NAME :%s INIT GRAPH NODES ------ >>> %+v", a.GetPath(), nodes))
	a.Graph.SetNodes(nodes)

	println(fmt.Sprintf("ACTION NAME :%s GET GRAPH NODES ------ >>> %+v", a.GetPath(), a.Graph.GetNodes()))
}

func (a *Action) GetGraph() interfaces.GraphInterface {
	return a.Graph
}

func (a *Action) OnRequest(method string, path string)  {

}

func (a *Action) Run(ctx context.Context, logger interfaces.LoggerInterface)  {

}

// */