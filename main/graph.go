package main

import (
	Gin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	ghvzExt "github.com/vortex14/gotyphoon/extensions/models/graphviz"
	GinExtensions "github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/resources/home"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func init()  {
	log.InitD()
}

func main()  {

	newG := (&ghvzExt.Graph{
		BaseGraph: &ghvzExt.BaseGraph{
			MetaInfo:  &label.MetaInfo{
				Name: "Server Graph",
			},
		},

	}).Init()


	_ = newG.AddSubGraph(&ghvzExt.GraphOptions{
		IsCluster: true,
		Label:     "pipeline group",
		Name:      "pipeline group",
		BackgroundColor: ghvzExt.COLORGold,
	}).AddNode(&ghvzExt.NodeOptions{
		Name:  "pipeline-№1",
		Shape: ghvzExt.SHAPEPipeline,
	}).AddNode(&ghvzExt.NodeOptions{
		EdgeOptions: &ghvzExt.EdgeOptions{
			NodeA: "pipeline-№1",
			Style: ghvzExt.StyleDotted,
		},
		Name: "pipeline-№2",
		Shape: ghvzExt.SHAPEPipeline,
	}).AddNode(&ghvzExt.NodeOptions{
		EdgeOptions: &ghvzExt.EdgeOptions{
			NodeA: "pipeline-№2",
			Style: ghvzExt.StyleDashed,
		},
		Name: "pipeline-№3",
		Shape: ghvzExt.SHAPEPipeline,
	}).AddNode(&ghvzExt.NodeOptions{
		EdgeOptions: &ghvzExt.EdgeOptions{
			NodeA: "pipeline-№3",
		},
		Name: "pipeline-№4",
		Shape: ghvzExt.SHAPEPipeline,
	})


	_ = newG.AddSubGraph(&ghvzExt.GraphOptions{
		IsCluster: true,
		Label:     "pipeline group 2",
		Name:      "pipeline group 2",
		BackgroundColor: ghvzExt.COLORAliceblue,
	}).AddNode(&ghvzExt.NodeOptions{
		Name:  "pipeline-№7",
		Shape: ghvzExt.SHAPEPipeline,
	}).AddNode(&ghvzExt.NodeOptions{
		EdgeOptions: &ghvzExt.EdgeOptions{
			NodeA: "pipeline-№7",
			Style: ghvzExt.StyleDotted,
		},
		Name: "pipeline-№8",
		Shape: ghvzExt.SHAPEPipeline,
	}).AddNode(&ghvzExt.NodeOptions{
		EdgeOptions: &ghvzExt.EdgeOptions{
			NodeA: "pipeline-№8",
			Style: ghvzExt.StyleDashed,
		},
		Name: "pipeline-№9",
		Shape: ghvzExt.SHAPEPipeline,
	})

	_ = newG.AddSubGraph(&ghvzExt.GraphOptions{
		IsCluster: true,
		Label:     "pipeline group 3",
		Name:      "pipeline group 3",
		BackgroundColor: ghvzExt.COLORBeige,
	}).AddNode(&ghvzExt.NodeOptions{
		Name:  "pipeline-№8",
		Shape: ghvzExt.SHAPEPipeline,
	}).AddNode(&ghvzExt.NodeOptions{
		EdgeOptions: &ghvzExt.EdgeOptions{
			NodeA: "pipeline-№8",
			Style: ghvzExt.StyleDotted,
		},
		Name: "pipeline-№9",
		Shape: ghvzExt.SHAPEPipeline,
	}).AddNode(&ghvzExt.NodeOptions{
		EdgeOptions: &ghvzExt.EdgeOptions{
			NodeA: "pipeline-№8",
			Style: ghvzExt.StyleDashed,
		},
		Name: "pipeline-№10",
		Shape: ghvzExt.SHAPEPipeline,
	})


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
					Path:        "graph",
					Name:        "image",
					Description: "Image data faker",
				},
				Methods: []string{interfaces.GET},
			},
			GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
				//wr := &bytes.Buffer{}
				//if err := g.Render(pipelines, graphviz.SVG, wr); err != nil {
				//	log.Fatal(err)
				//}
				//_, _ = ctx.Writer.Write(wr.Bytes())

				_, _ = ctx.Writer.Write(newG.Render("svg"))
			},
		}).AddAction(&GinExtensions.Action{
			Action: &forms.Action{
				MetaInfo: &label.MetaInfo{
					Name:        "Dot",
					Description: "Render Graph dot template",
					Path:        "graph-template",
				},

				Methods: []string{interfaces.GET},
			},

			GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {

				//var buf bytes.Buffer
				//
				//if err := g.Render(pipelines, "dot", &buf); err != nil {
				//	log.Fatal(err)
				//}
				//_, _ = ctx.Writer.Write(buf.Bytes())

				_, _ = ctx.Writer.Write(newG.Render("dot"))
			},
		})).Run()

	if errS != nil {
		logrus.Error(errS.Error())
	}
}
