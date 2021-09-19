package forms

import (
	Context "context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/elements/models/label"
	graphvizExt "github.com/vortex14/gotyphoon/extensions/models/graphviz"
	"github.com/vortex14/gotyphoon/interfaces"

	"github.com/vortex14/gotyphoon/log"
)

type PipelineGroup struct {
	*label.MetaInfo

	LambdaMap     map[string]interfaces.LambdaInterface
	PyLambdaMap   map[string]interfaces.LambdaInterface

	Stages        []interfaces.BasePipelineInterface
	Consumers     map[string]interfaces.ConsumerInterface

	graph         interfaces.GraphInterface
	LOG           interfaces.LoggerInterface

}

func (g *PipelineGroup) GetFirstPipelineName() string {
	firstStage := g.Stages[0]
	return firstStage.GetName()
}

func (g *PipelineGroup) Run(context Context.Context) {
	println("run pipeline group !")

	var failedFlow bool
	var mainContext Context.Context
	var middlewareContext Context.Context

	middlewareContext, mainContext = context, context

	mainContext = log.NewCtxValues(mainContext, log.D{"group": g.GetName()})

	for _, pipeline := range g.Stages {
		if failedFlow { break }
		logger := log.New(log.D{"pipeline": pipeline.GetName(), "group": g.GetName() })

		middlewareContext = log.NewCtx(mainContext, logger)
		var errStack error
		{
			var failed bool
			pipeline.RunMiddlewareStack(middlewareContext, func(middleware interfaces.MiddlewareInterface, err error) {
				errStack = err
				failed = true
				logger.Error("exit from middleware stack . Error: ", errStack.Error())
			}, func(returnedContext Context.Context) {
				middlewareContext = returnedContext
			})
			if failed { pipeline.Cancel(middlewareContext, logger, errStack); break }
		}

		mainContext = log.NewCtx(middlewareContext, logger)

		{
			pipeline.Run(mainContext, func(p interfaces.BasePipelineInterface, err error) {
				failedFlow = true
				errStack = err
				logger.Error("Exit from group. Error: ",err.Error(), p.GetName())
				p.Cancel(mainContext, logger, err)

			}, func(returnedResultPipelineContext Context.Context) {
				mainContext = returnedResultPipelineContext
			})
		}

	}
}

// /* ignore for building amd64-linux

func (g *PipelineGroup) InitGraph(parentNode string)  {
	groupGraph := g.graph.AddSubGraph(&interfaces.GraphOptions{
		IsCluster: true,
		Name: g.GetName(),
		Label: g.GetName(),
	})

	groupGraph.SetNodes(g.graph.GetNodes())

	g.LOG.Warning(">>>>>>>>>>>>>>>>>>",g.graph.GetNodes(), parentNode)
	var prevPipeline = parentNode
	for _, pipeline := range g.Stages {
		nodeOptions := &interfaces.NodeOptions{
			Name: pipeline.GetName(),
			Label: graphvizExt.FormatBottomSpace(pipeline.GetName()),
			Shape: graphvizExt.SHAPEPipeline,
			EdgeOptions: &interfaces.EdgeOptions{
				//NodeB:  a.handlerPath,
				ArrowS: 0.5,
			},
		}

		if len(prevPipeline) > 0 {
			nodeOptions.EdgeOptions.NodeA = prevPipeline
		}

		groupGraph.AddNode(nodeOptions)

		prevPipeline = pipeline.GetName()
	}
}

func (g *PipelineGroup) SetGraph(graph interfaces.GraphInterface)  {
	if g.graph != nil { return }
	g.LOG.Debug(fmt.Sprintf("SetGraph: %+v", graph.GetNodes()))
	g.graph = graph

}

func (g *PipelineGroup) SetGraphNodes(nodes map[string]interfaces.NodeInterface)  {
	logrus.Error(nodes)
	g.graph.SetNodes(nodes)
}

// */

func (g *PipelineGroup) SetLogger(logger interfaces.LoggerInterface)  {
	g.LOG = logger
}