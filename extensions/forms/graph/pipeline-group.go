package graph

import (
	"fmt"

	"github.com/vortex14/gotyphoon/elements/forms"
	graphvizExt "github.com/vortex14/gotyphoon/extensions/models/graphviz"
	"github.com/vortex14/gotyphoon/interfaces"
)

type PipelineGroup struct {
	*forms.PipelineGroup

	graph         interfaces.GraphInterface
}

func (g *PipelineGroup) InitGraph(parentNode string)  {
	groupGraph := g.graph.AddSubGraph(&interfaces.GraphOptions{
		IsCluster: true,
		Name: g.GetName(),
		Label: g.GetName(),
	})

	var prevPipeline = parentNode
	for _, pipeline := range g.Stages {
		nodeOptions := &interfaces.NodeOptions{
			Name: pipeline.GetName(),
			Label: graphvizExt.FormatBottomSpace(pipeline.GetName()),
			Shape: graphvizExt.SHAPEPipeline,
			EdgeOptions: &interfaces.EdgeOptions{
				//NodeB:  a.handlerPath,
				ArrowS: 0.5,
			},
		}

		if len(prevPipeline) > 0 {
			nodeOptions.EdgeOptions.NodeA = prevPipeline
		}

		groupGraph.AddNode(nodeOptions)

		prevPipeline = pipeline.GetName()
	}
}

func (g *PipelineGroup) SetGraph(logger interfaces.LoggerInterface, graph interfaces.GraphInterface)  {
	if g.graph != nil { return }
	logger.Debug(fmt.Sprintf("SetGraph: %+v", graph.GetNodes()))
	g.graph = graph
}

func (g *PipelineGroup) SetGraphNodes(nodes map[string]interfaces.NodeInterface)  {
	//logrus.Error(nodes)
	//g.graph.SetNodes(nodes)
}