package forms

import (
	"context"
	"github.com/sirupsen/logrus"
	Errors "github.com/vortex14/gotyphoon/errors"
	"sync"

	"github.com/vortex14/gotyphoon/interfaces"
)

type BasePipeline struct {
	Name        string
	Description string
	Required    bool
	//Task          *task.TyphoonTask
	//Project       interfaces.Project

	Context       context.Context
	//stageIndex    int32
	promise       sync.WaitGroup

	//inputCount    int64
	//inputByte     int64
	//InputMap      map[string]interface{}

	//outputCount   int64
	//outputByte    int64
	//OutputMap     map[string]interface{}


	//errorCount    int64
	//duration      time.Time


	//Handler     interfaces.HandlerInterface
	//Response    interfaces.ResponseInterface
	//
	//LOG         *logrus.Entry
	//Metrics     interfaces.MetricsInterface

	Middlewares []interfaces.MiddlewareInterface

	Callbacks   []interfaces.CallbackPipelineInterface
	//Consumers   map[string]interfaces.ConsumerInterface

	LambdaHandler func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context)
	CancelHandler func(ctx context.Context, err error)

	interfaces.BaseLabel
}

func (p *BasePipeline) NextStage()  {

}

func (p *BasePipeline) IsRequired() bool {
	return p.Required
}

func (p *BasePipeline) GetName() string {
	return p.Name
}

func (p *BasePipeline) GetDescription() string {
	return p.Description
}

func (p *BasePipeline) Await()  {
	p.promise.Wait()
}

func (p *BasePipeline) Run(ctx context.Context) (error, context.Context) {
	if p.LambdaHandler == nil {
		return Errors.LambdaRequired, nil
	}
	middlewareLogger := logrus.WithFields(logrus.Fields{
		"pipeline": p.GetName(),
	})
	return p.LambdaHandler(ctx, middlewareLogger)
}

func (p *BasePipeline) Cancel(err error) {
	p.CancelHandler(p.Context, err)
}

func (p *BasePipeline) RunMiddlewareStack(
	ctx context.Context,
	reject func(middleware interfaces.MiddlewareInterface,err error),
	next func(ctx context.Context),

	) {
	var failed bool
	var forceSkip bool
	var baseException error
	var middlewareContext context.Context

	middlewareContext = ctx
	for _, middleware := range p.Middlewares {
		if failed {break}
		if forceSkip {continue}

		middlewareLogger := logrus.WithFields(logrus.Fields{
			"middleware": middleware.GetName(), "pipeline": p.GetName(),
		})

		middleware.Pass(middlewareContext, middlewareLogger, func(err error) {
			if middleware.IsRequired() {baseException = err; err = Errors.MiddlewareRequired}

			switch err {
			case Errors.ForceSkipMiddlewares:
				forceSkip = true
				middlewareLogger.Warning(Errors.ForceSkipMiddlewares.Error())
			case Errors.MiddlewareRequired:
				reject(middleware, baseException)
				failed = true
			default:
				middlewareLogger.Warning(err.Error())
			}

		}, func(returnedMiddlewareContext context.Context) {
			middlewareContext = returnedMiddlewareContext
		})
		next(middlewareContext)
	}
}
