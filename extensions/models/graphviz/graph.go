package graphviz

//https:www.graphviz.org/pdf/dotguide.pdf

// /* ignore for building amd64-linux
//
//import (
//	"bytes"
//	"fmt"
//	"sync"
//
//	"github.com/goccy/go-graphviz"
//	"github.com/goccy/go-graphviz/cgraph"
//	"github.com/sirupsen/logrus"
//	. "github.com/vortex14/gotyphoon/elements/models/label"
//	. "github.com/vortex14/gotyphoon/elements/models/singleton"
//	Errors "github.com/vortex14/gotyphoon/errors"
//	"github.com/vortex14/gotyphoon/interfaces"
//	"github.com/vortex14/gotyphoon/log"
//)
//
//type Graph struct {
//	mu sync.Mutex
//	*MetaInfo
//	Singleton
//
//	LOG       *logrus.Entry
//	graph     *graphviz.Graphviz
//	Layout    string
//	Template  *cgraph.Graph
//
//	nodes     map[string] interfaces.NodeInterface
//	subGraphs map[string] interfaces.GraphInterface
//	edges     map[string] interfaces.EdgeInterface
//
//	Options *interfaces.GraphOptions
//}
//
//func (g *Graph) SetLayout(layout string) {
//	g.Layout = layout
//}
//
//func (g *Graph) Render(format string) []byte {
//	g.mu.Lock()
//	g.LOG.Debug("render server graph ...")
//	wr := &bytes.Buffer{}
//	if g.graph == nil {
//		return nil
//	}
//	if err := g.graph.Render(g.Template, GetExportFormat(format), wr); err != nil {
//		g.LOG.Debug("%s", err.Error())
//		return []byte{}
//	}
//	output := wr.Bytes()
//	g.LOG.Debug(fmt.Sprintf("output graph len: %d", len(output)))
//	g.mu.Unlock()
//	return output
//}
//
//
//func (g *Graph) GetEdgeName(a string, b string) string {
//	return fmt.Sprintf("%s->%s", a, b)
//}
//
//func (g *Graph) SetTemplate(template *cgraph.Graph)  {
//	g.Template = template
//}
//
//func (g *Graph) PostInit()  interfaces.GraphInterface {
//	g.LOG.Error("post init")
//	g.nodes = make(map[string] interfaces.NodeInterface)
//	g.edges = make(map[string] interfaces.EdgeInterface)
//
//	if g.Options == nil { return g}
//	var label string
//	var style string
//	var layout string
//	var colorFont string
//	var background string
//
//	if g.Options.IsCluster {
//		label = fmt.Sprintf("cluster-%s", g.Options.Label)
//	}
//
//	if len(g.Options.FontColor) == 0 { colorFont = COLORBlack } else {
//		colorFont = g.Options.FontColor
//	}
//
//	if len(g.Options.BackgroundColor) == 0 { background = COLORSnow } else {
//		background = g.Options.BackgroundColor
//	}
//
//	if len(g.Options.Style) == 0 { style = StyleSolid } else {
//		style = g.Options.Style
//	}
//
//	if len(g.Options.Layout) == 0 { layout = LAYOUTCirco } else {
//		layout = g.Options.Layout
//	}
//
//	g.LOG.Debug(
//		fmt.Sprintf("init subgraph. label: %s, bg: %s, fc: %s, style: %s, layout: %s",
//			label, background, colorFont, style, layout,
//		),
//	)
//
//	g.Template.
//		SubGraph(label, CFLAG).
//		SetLabel(label).
//		SetFontColor(colorFont).
//		SetStyle(GetStyle(style)).
//		SetBackgroundColor(background)
//
//	return g
//}
//
//func (g *Graph) AddNode(options *interfaces.NodeOptions ) interfaces.GraphInterface {
//	if options == nil { g.LOG.Error(Errors.GraphNodeOptionsNotFound.Error()); return nil}
//	if len(options.Label) == 0 { options.Label = options.Name }
//	node := (&Node{
//		MetaInfo: &MetaInfo{
//			Label: options.Label,
//			Name: options.Name,
//		},
//		LOG: g.LOG,
//		NodeOptions: options,
//	}).SetParent(g.Template).Init()
//
//	g.LOG.Warning(node)
//	g.LOG.Error(fmt.Sprintf(`
//------ Node name: %s
//Node options: %+v
//Defined nodes: %+v
//
//EdgeOptions: %+v
//
//`, options.Name, options, g.nodes, options.EdgeOptions))
//	if len(options.Style) > 0 { node.SetStyle(options.Style) }
//	if len(options.BackgroundColor) > 0 { node.SetColor(options.BackgroundColor) }
//	g.nodes[options.Name] = node
//
//	if options.EdgeOptions == nil { g.LOG.Warning("edge option not found") } else {
//		if nodeA, ok := g.nodes[options.EdgeOptions.NodeA]; ok {
//			edgeName := g.GetEdgeName(options.Name, options.EdgeOptions.NodeA)
//			g.LOG.Debug(fmt.Sprintf("init edge %s", edgeName), nodeA)
//			g.edges[edgeName] = (&Edge{
//				NodeB: node.Get(),
//				NodeA: nodeA.Get(),
//				EdgeOptions: options.EdgeOptions,
//				LOG: log.Patch(g.LOG, log.D{"edge": edgeName}),
//			}).SetGraph(g.Template).Init()
//		} else if nodeB, ok := g.nodes[options.EdgeOptions.NodeB]; ok {
//			edgeName := fmt.Sprintf("%s->%s", options.Name, options.EdgeOptions.NodeB)
//			g.LOG.Debug(fmt.Sprintf("init edge %s", edgeName), nodeB)
//			g.edges[edgeName] = (&Edge{
//				NodeB: nodeB.Get(),
//				NodeA: node.Get(),
//				EdgeOptions: options.EdgeOptions,
//				LOG: log.Patch(g.LOG, log.D{"edge": edgeName}),
//			}).SetGraph(g.Template).Init()
//		} else {
//			g.LOG.Warning(fmt.Sprintf("not found node A and B. %+v", options.EdgeOptions))
//		}
//	}
//
//	return g
//}
//
//func (g *Graph) AddSubGraph(options *interfaces.GraphOptions) interfaces.GraphInterface {
//	if options == nil { g.LOG.Error(Errors.GraphOptionsNotFound.Error()); return nil }
//	if len(options.Label) == 0 { g.LOG.Error(Errors.GraphOptionsLabelRequired.Error()); return nil }
//	if g.subGraphs == nil { g.subGraphs = make(map[string]interfaces.GraphInterface) }
//	if g.Template == nil { g.LOG.Error("template not found !") }
//	if len(options.Label) == 0 { g.LOG.Error(Errors.GraphNameNotFound.Error()); return nil}
//
//	subGraphLogger := log.Patch(g.LOG, log.D{"subgraph": options.Name})
//	subGraphLogger.Info("init sub graph ", options.Name)
//	g.subGraphs[options.Name] = (&Graph{
//		Template: g.Template.SubGraph(options.Name, CFLAG),
//		LOG: subGraphLogger,
//
//	}).PostInit()
//
//	return g.subGraphs[options.Name]
//}
//
//func (g *Graph) GetEdges() map[string]interfaces.EdgeInterface {
//	return g.edges
//}
//
//func (g *Graph) UpdateEdge(options *interfaces.EdgeOptions) interfaces.GraphInterface {
//	if options == nil { g.LOG.Error(Errors.GraphEdgeOptionsNotFound.Error()); return nil}
//	edgeName := g.GetEdgeName(options.NodeA, options.NodeB)
//	var selectedEdge interfaces.EdgeInterface
//	if edge, ok := g.edges[edgeName]; !ok {
//		g.LOG.Error(Errors.GraphEdgeNotFound.Error(), fmt.Sprintf(" edgeName: %s", edgeName))
//		return g
//	} else { selectedEdge = edge }
//
//	if len(options.Label)  > 0 { selectedEdge.SetLabel(options.Label)      }
//	if len(options.LabelH) > 0 { selectedEdge.SetHeadLabel(options.LabelH) }
//	if len(options.LabelT) > 0 { selectedEdge.SetTailLabel(options.LabelT) }
//
//	if options.ArrowS      > 0 { selectedEdge.SetArrowSize(options.ArrowS) }
//	if len(options.ArrowH) > 0 { selectedEdge.SetArrowHead(options.ArrowH) }
//	if len(options.ArrowT) > 0 { selectedEdge.SetArrowTail(options.ArrowT) }
//
//	if len(options.Color)  > 0 { selectedEdge.SetColor(options.Color)      }
//	if len(options.Style)  > 0 { selectedEdge.SetStyle(options.Style)      }
//
//	return g
//}
//
//func (s *Graph) BuildEdges(nodesA []string, nodesB []string) interfaces.GraphInterface {
//	for _, method := range nodesA {
//		nodeA := s.nodes[method]
//		for _, nodeName := range nodesB {
//			if len(nodeName) == 0 { continue }
//			nodeB := s.nodes[nodeName]
//			edgeName := s.GetEdgeName(method, nodeName)
//			s.edges[edgeName] = (&Edge{
//				NodeB: nodeB.Get(),
//				NodeA: nodeA.Get(),
//				EdgeOptions: &interfaces.EdgeOptions{},
//				LOG: log.Patch(s.LOG, log.D{"edge": edgeName}),
//			}).SetGraph(s.Template).Init()
//		}
//	}
//	return s
//}
//
//
//func (g *Graph) Init() interfaces.GraphInterface {
//	g.Construct(func() {
//
//		g.graph = graphviz.New()
//		g.graph.SetLayout(graphviz.Layout(g.Layout))
//		g.LOG = log.New(log.D{"graph": g.GetName()})
//		template, err := g.graph.Graph()
//		g.LOG.Error(g.Options)
//		if err != nil {
//			g.LOG.Error("err: %s", err.Error())
//			return
//		}
//		g.Template = template
//
//		g.PostInit()
//
//	})
//	return g
//}
//
//func (g *Graph) SetEdges(edges map[string]interfaces.EdgeInterface)  {
//	g.edges = edges
//}
//
//func (g *Graph) GetNodes() map[string]interfaces.NodeInterface {
//	return g.nodes
//}
//
//func (g *Graph) SetNodes(nodes map[string]interfaces.NodeInterface)  {
//	g.nodes = nodes
//}
//
//func (g *Graph) IsTemplate() bool {
//	status := false
//	if g.Template != nil { status = true }
//	return status
//}
//
//func (s *Graph) SetStyle(layout string)  {
//	s.Template.SetLayout(layout)
//}
////
// */