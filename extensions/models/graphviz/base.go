package graphviz

import (
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/sirupsen/logrus"

	. "github.com/vortex14/gotyphoon/elements/models/label"
	. "github.com/vortex14/gotyphoon/elements/models/singleton"
)

type BaseGraph struct {
	*MetaInfo
	Singleton

	template  *cgraph.Graph
	LOG       *logrus.Entry

	nodes     map[string] *Node
	subGraphs map[string] *SubGraph
	edges     map[string] *Edge


}

func (g *BaseGraph) SetTemplate(template *cgraph.Graph)  {
	g.template = template
}