package graphviz

import (
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/sirupsen/logrus"
	. "github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)

type Edge struct {
	Singleton
	*interfaces.EdgeOptions

	NodeA *cgraph.Node
	NodeB *cgraph.Node
	graph *cgraph.Graph
	LOG   *logrus.Entry

	edge  *cgraph.Edge
}


func (e *Edge) Init() *Edge {
	if utils.IsNill(e.NodeA, e.NodeB, e.EdgeOptions) { e.LOG.Debug( Errors.GraphEdgeContextBroken.Error()); return nil }
	e.Construct(func() {
		edge, _ := e.graph.CreateEdge(e.Name, e.NodeA, e.NodeB)
		e.edge = edge
		if len(e.EdgeOptions.Style) > 0 {
			e.edge.SetStyle(cgraph.EdgeStyle(e.EdgeOptions.Style))
		}

		if len(e.EdgeOptions.Label) > 0 {
			e.edge.SetLabel(e.EdgeOptions.Label)
		}

		if len(e.EdgeOptions.LabelH) > 0 {
			e.edge.SetHeadLabel(e.EdgeOptions.LabelH)
		}

		if len(e.EdgeOptions.LabelT) > 0 {
			e.edge.SetTailLabel(e.EdgeOptions.LabelT)
		}

		if e.EdgeOptions.ArrowS > 0 {
			e.edge.SetArrowSize(e.EdgeOptions.ArrowS)
		}

		if len(e.EdgeOptions.ArrowH) > 0 {
			e.edge.SetArrowHead(cgraph.ArrowType(e.EdgeOptions.ArrowH))
		}

		if len(e.EdgeOptions.ArrowT) > 0 {
			e.edge.SetArrowTail(cgraph.ArrowType(e.EdgeOptions.ArrowT))
		}

		if len(e.EdgeOptions.Color) > 0 {
			e.edge.SetColor(e.EdgeOptions.Color)
		}


	})
	return e
}

func (e *Edge) SetGraph(graph *cgraph.Graph) *Edge {
	e.graph = graph
	return e
}