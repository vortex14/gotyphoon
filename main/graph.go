package main

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	GinExtensions "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/resources/home"
	"log"
	"strconv"

	Gin "github.com/gin-gonic/gin"
	"github.com/goccy/go-graphviz"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	ARROWBox       = "box"
	ARROWVee       = "vee"
	ARROWNone      = "none"
	ARROWNormal    = "normal"

	COLORRed       = "red"
	COLORNavy      = "navy"
	COLORBlack     = "black"
	COLORGreen     = "green"
	COLORTomato    = "tomato"

	SHAPEBox3D     = "box3d"
	SHAPEFolder    = "folder"
	SHAPEPipeline  = "cylinder"
	SHAPEComponent = "component"


	LAYOUTDefault = "dot"
	LAYOUTNeato = "neato"
	LAYOUTLine= "sfdp"
	LAYOUTInvertedFlow = "twopi"

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

func CreateTestPipelineGroup(graph *cgraph.Graph, group string) *cgraph.Node {
	var prevStage *cgraph.Node

	for i:=0; i<5; i++  {
		cP := CreatePipeline(graph, fmt.Sprintf("Pipeline :%s: %s", group, strconv.Itoa(i)))

		if prevStage != nil {
			_, _ = graph.CreateEdge("stage 2", prevStage, cP)
		}
		prevStage = cP
	}

	return prevStage
}

func CreateBaseGraph()*cgraph.Graph  {
	tmpl := []byte(`digraph G {
subgraph cluster_0 {
color=red;
style=filled;
node [style=filled];
a0 -> a1;
label = "process #1";
}
}`)
	graph, _ := graphviz.ParseBytes(tmpl)
	return graph
}

func main()  {
	_ = []byte(`digraph G {

subgraph cluster_0 {
style=filled;
node [style=filled];
a0 -> a1;
label = "process #1";
}
}`)
	mainG := CreateBaseGraph()


	//mainG.SetBackgroundColor(COLORNavy)

	//gl, _ := graphviz.ParseBytes(tmpl)
	//gl.SetBackgroundColor(COLORTomato)
	//

	//for i := 0; i< gl.NumberSubGraph(); i ++ {
	//	println("subgraph",i, gl.SubGraph("cluster_1", 1).
	//		SetStyle("filled").
	//		//SetFontColor(COLORNavy).
	//		SetBackgroundColor(COLORBlack),
	//	)
	//}


	//if err != nil {
	//	return
	//}
	g := graphviz.New()

	pipelines, _ := g.Graph()
	pipeline := pipelines.SubGraph("test-1", 1).SetBackgroundColor(COLORBlack).SetStyle("filled")
	CreateTestPipelineGroup(pipeline, "g1")
	//pipelines.SetRootGraph(mainG)

	//mainG, _ := g.Graph()

	//mainG.SetBackgroundColor(COLORGreen)

	//neato
	//
	//{
	//	mainG.SetLabel("new group").SetLayout("sfdp")
	//
	//	sub := mainG.SubGraph("sub1", 1).
	//		SetStyle("filled").
	//		SetBackgroundColor(COLORTomato).
	//		SetBackground(COLORBlack).
	//		SetBackgroundColor(COLORTomato).
	//		SetFontColor(COLORRed)
	//
	//	CreateTestPipelineGroup(sub, "g1")
	//	//sub = sub.SetLabel("SUB1")
	//
	//	sub2 := sub.SubGraph("sub2", 2)
	//	CreateTestPipelineGroup(sub2, "g2")
	//}


	//newS := mainG.
	//println(newS.)
	//newS.SetColorScheme(COLORNavy)
	//newS.SetBackgroundColor(COLORBlack)
	//newS.SetStyle("filled")
	//CreateTestPipelineGroup(newS)
	
	//newS2 := mainG.SubGraph("test2", 1)
	//newS2.SetBackground(COLORTomato)
	//newS2.SetColorScheme(COLORNavy)
	//newS2.SetBackgroundColor(COLORBlack)
	//newS2.SetStyle("filled")
	//CreateTestPipelineGroup(newS2)

	//sub := mainG.SubGraph("sub1", 1)
	//sub.FirstSubGraph()
	println("COUNT: ",mainG.NumberSubGraph(), mainG)

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
	//
	//fmt.Println(11, buf.String())

	errS := (&GinExtensions.TyphoonGinServer{

		TyphoonServer: &forms.TyphoonServer{
			MetaInfo: &label.MetaInfo{
				Name:        "Graph Schema Generator",
				Description: "Generator Server Schema",
			},
			Port: 17668,
			IsDebug: true,
		},
	}).Init().AddResource(home.Constructor("/").
		AddAction(&GinExtensions.Action{
			Action: &forms.Action{
				MetaInfo: &label.MetaInfo{
					Name:        "image",
					Description: "Image data faker",
				},
				Path: "graph",
				Methods: []string{interfaces.GET},
			},
			GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
				wr := &bytes.Buffer{}
				if err := g.Render(pipelines, graphviz.SVG, wr); err != nil {
					log.Fatal(err)
				}
				_, _ = ctx.Writer.Write(wr.Bytes())
			},
		}).AddAction(&GinExtensions.Action{
			Action: &forms.Action{
				MetaInfo: &label.MetaInfo{
					Name:        "Dot",
					Description: "Render Graph dot template",
				},
				Path:        "graph-template",
				Methods: []string{interfaces.GET},
			},

			GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {

				var buf bytes.Buffer

				if err := g.Render(pipelines, "dot", &buf); err != nil {
					log.Fatal(err)
				}
				_, _ = ctx.Writer.Write(buf.Bytes())
			},
		})).Run()

	if errS != nil {
		logrus.Error(errS.Error())
	}
}
