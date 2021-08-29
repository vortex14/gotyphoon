package forms

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/interfaces"
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

func (p *PipelineGroup) GetLogger(name string) interfaces.LoggerInterface {
	return logrus.WithFields(logrus.Fields{
		"pipeline-group": p.GetName(),
		"pipeline": name,
	})
}

func (p *PipelineGroup) Run(ctx context.Context) {
	println("run pipeline group !")
	var pipelineContext context.Context
	var middlewareContext context.Context

	middlewareContext, pipelineContext = ctx, ctx

	for _, pipeline := range p.Stages {

		logger := p.GetLogger(pipeline.GetName())
		middlewareContext = interfaces.UpdateContext(ctx, interfaces.LOGGER, logger)

		{
			var failed bool
			pipeline.RunMiddlewareStack(middlewareContext, func(middleware interfaces.MiddlewareInterface, err error) {
				failed = true
				logger.Error("exit from middleware stack . Error: ", err.Error())
			}, func(returnedContext context.Context) {
				pipelineContext,middlewareContext = returnedContext, returnedContext

			})
			if failed { break }
		}





		{
			err, resultContext := pipeline.Run(pipelineContext)

			if err != nil {
				logger.Error("Exit from group. Error: ",err.Error(), pipeline.GetName())
				break
			} else if resultContext != nil {
				pipelineContext = resultContext
			}
		}

	}
}