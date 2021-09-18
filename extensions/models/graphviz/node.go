package graphviz
//
//import (
//	"fmt"
//	"github.com/fatih/color"
//	"github.com/goccy/go-graphviz/cgraph"
//	"github.com/sirupsen/logrus"
//	. "github.com/vortex14/gotyphoon/elements/models/label"
//	. "github.com/vortex14/gotyphoon/elements/models/singleton"
//	Errors "github.com/vortex14/gotyphoon/errors"
//	"github.com/vortex14/gotyphoon/interfaces"
//	"github.com/vortex14/gotyphoon/log"
//)
//
//
//type Node struct {
//	Singleton
//	*MetaInfo
//	*interfaces.NodeOptions
//
//	LOG      *logrus.Entry
//	node     *cgraph.Node
//	parent   *cgraph.Graph
//}
//
//func (n *Node) Init() interfaces.NodeInterface {
//	if n.parent == nil { color.Red(Errors.GraphParentNodeNotFound.Error()); return nil}
//
//	n.Construct(func() {
//		n.LOG = log.Patch(n.LOG, log.D{"node": n.GetName()})
//		n.LOG.Debug(fmt.Sprintf("init node %s shape: %s", n.GetName(), n.Shape))
//		node, _ := n.parent.CreateNode(n.GetName())
//		node.SetLabel(n.GetLabel())
//		node.SetShape(cgraph.Shape(n.Shape))
//		n.node = node
//	})
//	return n
//}
//
//func (n *Node) SetParent(parentGraph *cgraph.Graph) interfaces.NodeInterface {
//	n.parent = parentGraph
//	return n
//}
//
//func (n *Node) SetLabel(label string) interfaces.NodeInterface {
//	n.node.SetLabel(label)
//	return n
//}
//
//func (n *Node) Get() *cgraph.Node {
//	return n.node
//}
//
//func (n *Node) SetStyle(style string)  {
//	n.node.SetStyle(cgraph.NodeStyle(style))
//}
//
//func (n *Node) SetColor(color string)  {
//	n.node.SetFillColor(color)
//}