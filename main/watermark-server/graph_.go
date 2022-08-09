package main

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/resources/home"
	"log"
	"strconv"

	Gin "github.com/gin-gonic/gin"
	"github.com/goccy/go-graphviz"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	ARROWBox    = "box"
	ARROWVee    = "vee"
	ARROWNone   = "none"
	ARROWNormal = "normal"

	COLORRed    = "red"
	COLORNavy   = "navy"
	COLORBlack  = "black"
	COLORGreen  = "green"
	COLORTomato = "tomato"

	SHAPEBox3D     = "box3d"
	SHAPEFolder    = "folder"
	SHAPEPipeline  = "cylinder"
	SHAPEComponent = "component"
)

//digraph G {
//layers="local:pvt:test:new:ofc";
//
//node1  [layer="pvt"];
//node2  [layer="all"];
//node3  [layer="pvt:ofc"];		/* pvt, test, new, and ofc */
//node2 -> node3  [layer="pvt:all"];	/* same as pvt:ofc */
//node2 -> node4 [layer=3];		/* same as test */
//}

func CreatePipeline(graph *cgraph.Graph, name string) *cgraph.Node {
	p, _ := graph.CreateNode(name)
	return p.SetShape(SHAPEPipeline)
}

func CreateTestPipelineGroup(graph *cgraph.Graph) *cgraph.Node {
	var prevStage *cgraph.Node

	for i := 0; i < 5; i++ {
		cP := CreatePipeline(graph, fmt.Sprintf("Pipeline :: %s", strconv.Itoa(i)))

		if prevStage != nil {
			_, _ = graph.CreateEdge("stage 2", prevStage, cP)
		}
		prevStage = cP
	}

	return prevStage
}

func main2() {
	tmpl := []byte(`digraph G {

subgraph cluster_0 {
style=filled;
node [style=filled,color=white];
a0 -> a1 -> a2 -> a3;
label = "process #1";
}

subgraph cluster_1 {
node [style=filled];
b0 -> b1 -> b2 -> b3;
label = "process #2";
color=blue
}
start -> a0;
start -> b0;
a1 -> b3;
b2 -> a3;
a3 -> a0;
a3 -> end;
b3 -> end;

start [shape=Mdiamond];
end [shape=Msquare];
}`)

	gl, err := graphviz.ParseBytes(tmpl)
	gl.SetBackgroundColor(COLORTomato)

	for i := 0; i < gl.NumberSubGraph(); i++ {
		println("subgraph", i, gl.SubGraph("cluster_0", 1).
			SetStyle("filled").
			SetFontColor(COLORNavy).
			SetBackgroundColor(COLORBlack),
		)
	}
	//if err != nil {
	//	return
	//}
	g := graphviz.New()

	mainG, _ := g.Graph()

	newS := mainG.SubGraph("test", -1)
	newS.SetBackground(COLORTomato)
	newS.SetColorScheme(COLORNavy)
	newS.SetBackgroundColor(COLORBlack)
	newS.SetStyle("filled")
	CreateTestPipelineGroup(newS)

	//newS2 := mainG.SubGraph("test2", 1)
	//newS2.SetBackground(COLORTomato)
	//newS2.SetColorScheme(COLORNavy)
	//newS2.SetBackgroundColor(COLORBlack)
	//newS2.SetStyle("filled")
	//CreateTestPipelineGroup(newS2)

	//sub := mainG.SubGraph("sub1", 1)
	//sub.FirstSubGraph()
	println(mainG.NumberSubGraph(), newS)

	//
	//CreateTestPipelineGroup(sub)
	//
	//sub2 := mainG.SubGraph("sub2", 2)
	//CreateTestPipelineGroup(sub2)

	//newSubGraph := sub.SubGraph("test", 1)
	//
	//
	//graph, err := graphviz.ParseBytes(tmpl)
	//
	//newSubGraph = newSubGraph.SetPage(20)
	//
	////println(graph.SubGraph("test", 1).)
	//
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer func() {
	//	if err := newSubGraph.Close(); err != nil {
	//		log.Fatal(err)
	//	}
	//	g.Close()
	//}()
	//
	//n, err := newSubGraph.CreateNode("n")
	//n.SetFillColor(COLORBlack)
	//n.SetArea(200)
	//n.SetLayer("layer1")
	//
	//n.SetColor(COLORRed)
	//
	//n.SetLabel("test label 11")
	//n.SetComment("test comment")
	//n.SetGroup("layer2")
	//
	////println("LEN NODES ", g.NumberSubGraph())
	//
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//m, err := newSubGraph.CreateNode("p1")
	//m.SetLayer("layer2")
	//m.SetArea(200)
	//m.SetGroup("layer2")
	//m.SetFillColor(COLORTomato)
	//m.SetLayer("test-layer")
	//m.SetShape(SHAPEPipeline)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//p, _ := newSubGraph.CreateNode("p2")
	//
	//p.SetShape(SHAPEPipeline)
	//e, err := newSubGraph.CreateEdge("----edge---", n, m)
	//e.SetLayer("layer1:layer2")
	//
	//e.SetColor("red")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, _ = newSubGraph.CreateEdge("stage 2", m, p)
	//
	//e.SetLabel("label for line ")
	var buf bytes.Buffer

	if err := g.Render(gl, "dot", &buf); err != nil {
		log.Fatal(err)
	}
	fmt.Println(11, buf.String())

	err = (&gin.TyphoonGinServer{
		TyphoonServer: &forms.TyphoonServer{
			MetaInfo: &label.MetaInfo{
				Name:        "Graph Schema Generator",
				Description: "Generator Server Schema",
			},
			Port:    17668,
			IsDebug: true,
		},
	}).Init().InitLogger().AddResource(home.Constructor("/").AddAction(&gin.Action{
		Action: &forms.Action{
			MetaInfo: &label.MetaInfo{
				Name:        "image",
				Description: "Image data faker",
				Path:        "graph",
			},
			Methods: []string{interfaces.GET},
		},
		GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
			wr := &bytes.Buffer{}
			if err := g.Render(gl, graphviz.SVG, wr); err != nil {
				log.Fatal(err)
			}
			_, _ = ctx.Writer.Write(wr.Bytes())
		},
	})).Run()

	if err != nil {
		logrus.Error(err.Error())
	}
}
