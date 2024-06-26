package forms

import (
	Context "context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/log"
	"golang.org/x/sync/semaphore"
	"sync"

	// /* ignore for building amd64-linux
	//	"fmt"
	//	"github.com/sirupsen/logrus"
	//	graphvizExt "github.com/vortex14/gotyphoon/extensions/models/graphviz"
	//*/
	"github.com/vortex14/gotyphoon/interfaces"
)

//const name =

type PipelineGroup struct {
	*label.MetaInfo

	sem *semaphore.Weighted

	LambdaMap   map[string]interfaces.LambdaInterface
	PyLambdaMap map[string]interfaces.LambdaInterface

	Options *Options
	ctx     Context.Context
	syncCtx sync.Once

	Stages    []interfaces.BasePipelineInterface
	Consumers map[string]interfaces.ConsumerInterface
	// /* ignore for building amd64-linux
	//	graph interfaces.GraphInterface
	//*/
	LOG interfaces.LoggerInterface
}

func (g *PipelineGroup) GetFirstPipelineName() string {
	firstStage := g.Stages[0]
	return firstStage.GetName()
}

func (g *PipelineGroup) initSemaphore() bool {
	status := true

	if g.Options != nil && g.sem == nil && g.Options.MaxConcurrent > 0 {
		g.sem = semaphore.NewWeighted(g.Options.MaxConcurrent)
	}

	if g.sem != nil && !g.sem.TryAcquire(1) {
		status = false
	}

	return status
}

func (g *PipelineGroup) initCtx() {
	g.syncCtx.Do(func() {
		g.ctx = Context.Background()
		if g.Options == nil {
			g.Options = &Options{}
		}
	})
}

func (g *PipelineGroup) Run(context Context.Context) error {

	var forceSkip bool
	var failedFlow bool
	var mainContext Context.Context

	mainContext = context
	logState := log.D{"group": g.GetName(), "call_id": uuid.New().String()}
	mainContext = log.PatchCtx(mainContext, logState)

	semStatus := g.initSemaphore()
	g.initCtx()

	if !semStatus {
		return Errors.PipelineCrowded
	}

	var errStack error
	for index, pipeline := range g.Stages {
		if failedFlow || forceSkip {
			break
		}
		logState["pipeline"] = pipeline.GetName()
		mainContext = log.PatchCtx(mainContext, logState)

		_, logger := log.Get(mainContext)

		if skipFlag, numberStage := GetGOTOCtx(mainContext); skipFlag && numberStage > index+1 {
			continue
		}

		if !g.Options.NotSharedContext {
			pipeline.SetSharedCtx(&g.ctx)
		}

		pipeline.Run(mainContext, func(canceledCtx Context.Context, p interfaces.BasePipelineInterface, err error) {
			switch {
			case errors.Is(err, Errors.ForceSkipPipelines):
				forceSkip = true
				logger.Warn(Errors.ForceSkipPipelines.Error())
			default:
				errStack = err
				logger.Error(fmt.Sprintf("Pipeline name: %s ; Exit from group. Error: %s", p.GetName(), err.Error()))
				failedFlow = true
				pipeline.Cancel(canceledCtx, logger, errStack)
			}

		}, func(returnedResultPipelineContext Context.Context) {
			errStack = nil
			mainContext = returnedResultPipelineContext
		})

	}

	if g.sem != nil {
		g.sem.Release(1)
	}

	return errStack
}

// /* ignore for building amd64-linux
//
//func (g *PipelineGroup) InitGraph(parentNode string) {
//	groupGraph := g.graph.AddSubGraph(&interfaces.GraphOptions{
//		IsCluster: true,
//		Name:      g.GetName(),
//		Label:     g.GetName(),
//	})
//
//	groupGraph.SetNodes(g.graph.GetNodes())
//
//	g.LOG.Warning(">>>>>>>>>>>>>>>>>>", g.graph.GetNodes(), parentNode)
//	var prevPipeline = parentNode
//	for _, pipeline := range g.Stages {
//		nodeOptions := &interfaces.NodeOptions{
//			Name:  pipeline.GetName(),
//			Label: graphvizExt.FormatBottomSpace(pipeline.GetName()),
//			Shape: graphvizExt.SHAPEPipeline,
//			EdgeOptions: &interfaces.EdgeOptions{
//				//NodeB:  a.handlerPath,
//				ArrowS: 0.5,
//			},
//		}
//
//		if len(prevPipeline) > 0 {
//			nodeOptions.EdgeOptions.NodeA = prevPipeline
//		}
//
//		groupGraph.AddNode(nodeOptions)
//
//		prevPipeline = pipeline.GetName()
//	}
//}
//
//func (g *PipelineGroup) SetGraph(graph interfaces.GraphInterface) {
//	if g.graph != nil {
//		return
//	}
//	g.LOG.Debug(fmt.Sprintf("SetGraph: %+v", graph.GetNodes()))
//	g.graph = graph
//
//}
//
//func (g *PipelineGroup) SetGraphNodes(nodes map[string]interfaces.NodeInterface) {
//	logrus.Error(nodes)
//	g.graph.SetNodes(nodes)
//}
//
// */

func (g *PipelineGroup) SetLogger(logger interfaces.LoggerInterface) {
	g.LOG = logger
}
