package graphviz

import (
	"fmt"
	"github.com/vortex14/gotyphoon/interfaces"

	"github.com/fatih/color"
	"github.com/goccy/go-graphviz/cgraph"

	. "github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/log"
)

type SubGraph struct {
	*BaseGraph

	parent  *cgraph.Graph
	graph   *cgraph.Graph

	Options *interfaces.GraphOptions
}

func (s *SubGraph) SetParent(parent *cgraph.Graph) *SubGraph {
	s.parent = parent
	return s
}

func (s *SubGraph) AddSubGraph(options *interfaces.GraphOptions) interfaces.GraphInterface {
	s.LOG.Error("not implement")
	return s
}

func (s *SubGraph) Render(format string) []byte {
	s.LOG.Error("not implement")
	return nil
}

func (s *SubGraph) UpdateEdge(options *interfaces.EdgeOptions) interfaces.GraphInterface {
	if options == nil { s.LOG.Error(Errors.GraphEdgeOptionsNotFound.Error()); return nil}

	return s
}

func (s *SubGraph) AddNode(options *interfaces.NodeOptions ) interfaces.GraphInterface {
	if options == nil { s.LOG.Error(Errors.GraphNodeOptionsNotFound.Error()); return nil}

	node := (&Node{
		MetaInfo: &MetaInfo{
			Name: options.Name,
		},
		LOG: s.LOG,
		NodeOptions: options,
	}).SetParent(s.template).Init()

	s.nodes[options.Name] = node

	if options.EdgeOptions == nil { s.LOG.Warning("edge option not found") } else {
		if nodeA, ok := s.nodes[options.EdgeOptions.NodeA]; ok {
			edgeName := fmt.Sprintf("%s->%s", options.EdgeOptions.NodeA, options.Name)
			s.LOG.Debug("init edge %s", edgeName)
			s.edges[edgeName] = (&Edge{
				NodeB: node.node,
				NodeA: nodeA.node,
				EdgeOptions: options.EdgeOptions,
				LOG: log.Patch(s.LOG, log.D{"edge": edgeName}),
			}).SetGraph(s.parent).Init()
		} else if nodeB, ok := s.nodes[options.EdgeOptions.NodeB]; ok {
			edgeName := fmt.Sprintf("%s->%s", options.EdgeOptions.NodeB, options.Name)
			s.LOG.Debug("init edge %s", edgeName)
			s.edges[edgeName] = (&Edge{
				NodeB: nodeB.node,
				NodeA: node.node,
				EdgeOptions: options.EdgeOptions,
				LOG: log.Patch(s.LOG, log.D{"edge": edgeName}),
			}).SetGraph(s.parent).Init()
		} else {
			s.LOG.Warning("not found node A and B")
		}
	}

	return s
}

func (s *SubGraph) Init() interfaces.GraphInterface {
	if s.parent == nil { color.Red(Errors.GraphMainGraphNotFound.Error()) }
	s.Construct(func() {

		s.nodes = make(map[string] *Node)
		s.edges = make(map[string] *Edge)

		var label string
		var style string
		var colorFont string
		var background string

		if s.Options.IsCluster {
			label = fmt.Sprintf("cluster-%s", s.Options.Label)
		}

		if len(s.Options.FontColor) == 0 { colorFont = COLORBlack } else {
			colorFont = s.Options.FontColor
		}

		if len(s.Options.BackgroundColor) == 0 { background = COLORSnow } else {
			background = s.Options.BackgroundColor
		}

		if len(s.Options.Style) == 0 { style = StyleSolid } else {
			style = s.Options.Style
		}

		s.LOG.Debug(
			fmt.Sprintf("init subgraph. label: %s, bg: %s, fc: %s, style: %s",
				label, background, colorFont, style,
			),
		)

		s.template = s.parent.
			SubGraph(label, CFLAG).
			SetLabel(label).
			SetFontColor(colorFont).
			SetStyle(GetStyle(style)).
			SetBackgroundColor(background)
	})

	return s
}