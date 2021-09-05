package forms

import (
	Context "context"

	"github.com/vortex14/gotyphoon/elements/models/awaitable"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/utils"

	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type BasePipeline struct {
	*label.MetaInfo
	*awaitable.Object

	//stageIndex    int32

	//inputCount    int64
	//inputByte     int64
	//InputMap      map[string]interface{}

	//outputCount   int64
	//outputByte    int64
	//OutputMap     map[string]interface{}

	//errorCount    int64
	//duration      time.Time

	//Metrics     interfaces.MetricsInterface

	Middlewares []interfaces.MiddlewareInterface

	Callbacks   []interfaces.CallbackPipelineInterface
	//Consumers   map[string]interfaces.ConsumerInterface

	Fn func(ctx Context.Context, logger interfaces.LoggerInterface) (error, Context.Context)
	Cn func(ctx Context.Context, logger interfaces.LoggerInterface, err error)
}

func (p *BasePipeline) Run(
	context Context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx Context.Context),
	) {
	if utils.IsNill(p.Fn) { reject(p,Errors.LambdaRequired); return}
	var logCtx interfaces.LoggerInterface
	if ok, logger := log.Get(context); !ok { reject(p, Errors.CtxLogFailed); return } else { logCtx = logger }
	err, newContext := p.Fn(context, logCtx)
	if utils.NotNill(err) { reject(p, err); return }

	next(newContext)
}

func (p *BasePipeline) Cancel(ctx Context.Context, logger interfaces.LoggerInterface, err error) {
	if utils.IsNill(p.Cn) { return }
	p.Cn(ctx, logger, err)
}

func (p *BasePipeline) RunMiddlewareStack(
	context Context.Context,
	reject func(middleware interfaces.MiddlewareInterface,err error),
	next func(ctx Context.Context),

	) {
	var failed bool
	var forceSkip bool
	var baseException error
	var middlewareContext Context.Context

	middlewareContext = context
	for _, middleware := range p.Middlewares {
		if failed || forceSkip { break }

		logger := log.New(log.D{"middleware": middleware.GetName(), "pipeline": p.GetName()})

		middleware.Pass(middlewareContext, logger, func(err error) {
			if middleware.IsRequired() {baseException = err; err = Errors.MiddlewareRequired}
			switch err {
			case Errors.ForceSkipMiddlewares:
				forceSkip = true
				logger.Warning(Errors.ForceSkipMiddlewares.Error())
			case Errors.MiddlewareRequired:
				reject(middleware, baseException)
				failed = true
			default:
				logger.Warning(err.Error())
			}

		}, func(returnedMiddlewareContext Context.Context) {
			middlewareContext = returnedMiddlewareContext
		})
		next(middlewareContext)
	}
}
