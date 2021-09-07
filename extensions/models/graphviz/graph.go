package graphviz

import (
	"bytes"
	"fmt"

	"github.com/fatih/color"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"

	. "github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/log"
)

type GraphOptions struct {
	Name            string
	IsCluster       bool
	FontColor       string
	Label           string
	BackgroundColor string
	PrefixNodeName  string
	Style           string
}


type Graph struct {
	*BaseGraph

	template    *cgraph.Graph
	graph       *graphviz.Graphviz
}

func (g *Graph) Init() *Graph {
	g.Construct(func() {
		g.graph = graphviz.New()
		g.LOG = log.New(log.D{"graph": g.GetName()})
		template, err := g.graph.Graph()
		if err != nil {
			g.LOG.Error("err: %s", err.Error())
			return
		}
		g.template = template
		g.LOG.Debug(fmt.Sprintf("init graph %s", g.GetName()))
	})
	return g
}



func (g *Graph) AddSubGraph(options *GraphOptions) *SubGraph {
	if options == nil { g.LOG.Error(Errors.GraphOptionsNotFound.Error()); return nil }
	if len(options.Label) == 0 { g.LOG.Error(Errors.GraphOptionsLabelRequired.Error()); return nil }
	if g.subGraphs == nil { g.subGraphs = make(map[string]*SubGraph) }

	if len(options.Label) == 0 { g.LOG.Error(Errors.GraphNameNotFound.Error()); return nil}

	subGraphLogger := log.Patch(g.LOG, log.D{"subgraph": options.Name})
	g.subGraphs[options.Name] = (&SubGraph{
		Options: options,
		BaseGraph: &BaseGraph{
			MetaInfo: &MetaInfo{Name: options.Name},
			LOG: subGraphLogger,
		},
	}).SetParent(g.template).Init()

	return g.subGraphs[options.Name]
}

func (g *Graph) Render(format string) []byte {
	wr := &bytes.Buffer{}
	if err := g.graph.Render(g.template, GetExportFormat(format), wr); err != nil {
		color.Red("%s", err.Error())
		return nil
	}
	return wr.Bytes()
}