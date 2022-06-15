package tests

import (
	"bytes"
	"log"
	"testing"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"

	. "github.com/smartystreets/goconvey/convey"
)

func Export(graph *graphviz.Graphviz, template *cgraph.Graph) []byte {
	wr := &bytes.Buffer{}
	if err := graph.Render(template, graphviz.XDOT, wr); err != nil {
		log.Fatalln(err.Error())
		return []byte{}
	}
	return wr.Bytes()
}

func TestExportGraph(t *testing.T) {

	Convey("init graphviz", t, func() {
		graph := graphviz.New()
		graph.SetLayout(graphviz.CIRCO)
		template, err := graph.Graph()
		var subGraph *cgraph.Graph
		So(err, ShouldBeNil)
		Convey("create new subgraph", func() {
			subGraph = template.SubGraph("cluster-1", 1)
			So(subGraph, ShouldNotBeEmpty)

			Convey("create new node a and b", func() {
				_, errNodeA := subGraph.CreateNode("node-A---")
				_, errNodeB := subGraph.CreateNode("node-B----")
				So(errNodeA, ShouldBeNil)
				So(errNodeB, ShouldBeNil)

				Convey("export to dot format", func() {
					o := Export(graph, template)
					So(o, ShouldNotBeEmpty)

					result := `digraph "" {
	graph [bb="0,0,195.96,92"];
	node [label="\N"];
	subgraph "cluster-1" {
		"node-A---"		 [height=0.5,
			pos="146.32,74",
			width=1.3789];
		"node-B----"		 [height=0.5,
			pos="52.323,18",
			width=1.4534];
	}
}`
					So(string(o), ShouldContainSubstring, result)
				})
			})

		})

	})

}
