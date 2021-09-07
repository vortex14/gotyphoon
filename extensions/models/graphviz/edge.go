package graphviz

import (
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/sirupsen/logrus"

	. "github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/utils"
)

type Edge struct {
	Singleton
	*EdgeOptions

	NodeA *cgraph.Node
	NodeB *cgraph.Node
	graph *cgraph.Graph
	LOG   *logrus.Entry

	edge  *cgraph.Edge
}

type EdgeOptions struct {
	Name string
	Arrow string
	Style string
	NodeA string
	NodeB string
}

func (e *Edge) Init() *Edge {
	if utils.IsNill(e.NodeA, e.NodeB, e.EdgeOptions) { e.LOG.Debug( Errors.GraphEdgeContextBroken.Error()); return nil }
	e.Construct(func() {
		edge, _ := e.graph.CreateEdge(e.Name, e.NodeA, e.NodeB)
		e.edge = edge
		if len(e.EdgeOptions.Style) > 0 {
			e.edge.SetStyle(cgraph.EdgeStyle(e.EdgeOptions.Style))
		}
	})
	return e
}

func (e *Edge) SetGraph(graph *cgraph.Graph) *Edge {
	e.graph = graph
	return e
}