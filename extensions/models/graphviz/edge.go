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

func (e *Edge) Init() interfaces.EdgeInterface {
	if utils.IsNill(e.NodeA, e.NodeB, e.EdgeOptions, e.graph) { e.LOG.Debug( Errors.GraphEdgeContextBroken.Error()); return nil }
	e.Construct(func() {
		e.LOG.Debug("init edge ", e.graph.Name(), "!!!!!!!!")
		edge, _ := e.graph.CreateEdge(e.Name, e.NodeA, e.NodeB)

		e.edge = edge

		if len(e.EdgeOptions.Label)  > 0 { e.SetLabel(e.EdgeOptions.Label)      }
		if len(e.EdgeOptions.LabelH) > 0 { e.SetHeadLabel(e.EdgeOptions.LabelH) }
		if len(e.EdgeOptions.LabelT) > 0 { e.SetTailLabel(e.EdgeOptions.LabelT) }

		if e.EdgeOptions.ArrowS      > 0 { e.SetArrowSize(e.EdgeOptions.ArrowS) }
		if len(e.EdgeOptions.ArrowH) > 0 { e.SetArrowHead(e.EdgeOptions.ArrowH) }
		if len(e.EdgeOptions.ArrowT) > 0 { e.SetArrowTail(e.EdgeOptions.ArrowT) }

		if len(e.EdgeOptions.Style)  > 0 { e.SetStyle(e.EdgeOptions.Style)      }
		if len(e.EdgeOptions.Color)  > 0 { e.SetColor(e.EdgeOptions.Color)      }
	})
	return e
}

func (e *Edge) SetGraph(graph *cgraph.Graph) interfaces.EdgeInterface {
	e.graph = graph
	e.LOG.Debug("SetGraph", e.Name, graph)
	return e
}

func (e *Edge) SetLabel(label string)  {
	e.edge.SetLabel(label)
}

func (e *Edge) SetStyle(style string)  {
	e.edge.SetStyle(cgraph.EdgeStyle(style))
}

func (e *Edge) SetColor(color string)  {
	e.edge.SetColor(color)
}

func (e *Edge) SetArrowSize(size float64)  {
	e.edge.SetArrowSize(size)
}

func (e *Edge) SetHeadLabel(label string)  {
	e.edge.SetHeadLabel(label)
}

func (e *Edge) SetArrowHead(head string)  {
	e.edge.SetArrowHead(cgraph.ArrowType(head))
}

func (e *Edge) SetTailLabel(label string)  {
	e.edge.SetTailLabel(label)
}

func (e *Edge) SetArrowTail(tail string)  {
	e.edge.SetArrowTail(cgraph.ArrowType(tail))
}