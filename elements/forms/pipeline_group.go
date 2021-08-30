package forms

import (
	Context "context"
	"github.com/vortex14/gotyphoon/interfaces"
	"golang.org/x/net/context"

	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/log"
)

type PipelineGroup struct {
	interfaces.BaseLabel
	//Name string
	//Description string
	//Required bool

	//errorCount    int64
	//duration      time.Time
	//timeLife      time.Time

	LambdaMap     map[string]interfaces.LambdaInterface
	PyLambdaMap   map[string]interfaces.LambdaInterface

	Stages        []interfaces.BasePipelineInterface
	Consumers     map[string]interfaces.ConsumerInterface

}


func (g *PipelineGroup) Run(context context.Context) {
	println("run pipeline group !")

	var failedFlow bool
	var mainContext Context.Context
	var middlewareContext Context.Context

	middlewareContext, mainContext = context, context

	mainContext = ctx.UpdateContext(mainContext, interfaces.LOGGER, log.GetContext(log.D{"group": g.GetName()}))

	for _, pipeline := range g.Stages {
		if failedFlow { break }
		logger := log.GetContext(log.D{"pipeline": pipeline.GetName(), "group": g.GetName() })

		middlewareContext = ctx.UpdateContext(mainContext, interfaces.LOGGER, logger)

		{
			var failed bool
			pipeline.RunMiddlewareStack(middlewareContext, func(middleware interfaces.MiddlewareInterface, err error) {
				failed = true
				logger.Error("exit from middleware stack . Error: ", err.Error())
			}, func(returnedContext Context.Context) {
				middlewareContext = returnedContext
			})
			if failed { break }
		}

		mainContext = ctx.UpdateContext(middlewareContext, interfaces.LOGGER, logger)

		{
			pipeline.Run(mainContext, func(pipeline interfaces.BasePipelineInterface, err error) {
				failedFlow = true
				logger.Error("Exit from group. Error: ",err.Error(), pipeline.GetName())

			}, func(returnedResultPipelineContext Context.Context) {
				mainContext = returnedResultPipelineContext
			})


		}

	}
}