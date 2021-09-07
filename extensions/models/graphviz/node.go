package graphviz

import (
	"github.com/fatih/color"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/sirupsen/logrus"
	. "github.com/vortex14/gotyphoon/elements/models/label"
	. "github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/log"
)


type NodeOptions struct {
	*EdgeOptions
	Name string
	Shape string
}

type Node struct {
	Singleton
	*MetaInfo
	*NodeOptions

	LOG      *logrus.Entry
	node     *cgraph.Node
	parent   *cgraph.Graph
}

func (n *Node) Init() *Node {
	if n.parent == nil { color.Red(Errors.GraphParentNodeNotFound.Error()); return nil}
	n.Construct(func() {
		n.LOG = log.Patch(n.LOG, log.D{"node": n.GetName()})
		n.LOG.Debug("init node")
		node, _ := n.parent.CreateNode(n.GetName())
		n.node = node
		node.SetShape(cgraph.Shape(n.Shape))
	})
	return n
}

func (n *Node) SetParent(parentGraph *cgraph.Graph) *Node {
	n.parent = parentGraph
	return n
}

