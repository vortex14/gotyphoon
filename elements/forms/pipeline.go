package forms

import (
	Context "context"
	"errors"
	"fmt"
	"github.com/avast/retry-go/v4"
	"time"

	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/utils"

	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

const (
	RetryCount     = "retry_count"
	PanicException = "PANIC"
)

type RetryOptions struct {
	Sleep time.Duration

	MaxCount int

	Required            bool
	OnlyRetryExceptions bool

	RetryExceptions    []error
	CriticalExceptions []error
}

type Options struct {
	Retry RetryOptions
}

func GetDefaultRetryOptions() *Options {
	return &Options{Retry: RetryOptions{
		MaxCount: 7,
	}}
}

func GetNotRetribleOptions() *Options {
	return &Options{
		Retry: RetryOptions{
			MaxCount: 1,
		},
	}
}

func GetCustomRetryOptions(count int) *Options {
	return &Options{
		Retry: RetryOptions{
			MaxCount: count,
		},
	}
}

type BasePipeline struct {
	*label.MetaInfo
	*awaitabler.Object

	Options        *Options
	NotIgnorePanic bool

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

	Callbacks []interfaces.CallbackPipelineInterface
	//Consumers   map[string]interfaces.ConsumerInterface

	Fn func(ctx Context.Context, logger interfaces.LoggerInterface) (error, Context.Context)
	Cn func(ctx Context.Context, logger interfaces.LoggerInterface, err error)
}

func (p *BasePipeline) run(
	context Context.Context,
	logCtx interfaces.LoggerInterface,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx Context.Context),
) {

	if utils.IsNill(p.Fn) {
		reject(p, Errors.LambdaRequired)
		p.Cancel(context, logCtx, Errors.LambdaRequired)
		return
	}

	err, newContext := p.Fn(context, logCtx)
	if utils.NotNill(err) {
		reject(p, err)
		p.Cancel(context, logCtx, err)
		return
	}

	next(newContext)
}

func (p *BasePipeline) SafeRun(run func() error, catch func(err error)) {
	if p.Options == nil {
		p.Options = GetDefaultRetryOptions()
	}
	defer func() {
		if !p.NotIgnorePanic {
			if r := recover(); r != nil {
				panicE := errors.New(fmt.Sprintf("%s: %s", PanicException, r))
				catch(panicE)
			}
		}

	}()
	err := retry.Do(func() error {
		return run()
	},
		retry.Attempts(uint(p.Options.Retry.MaxCount)),
		retry.RetryIf(func(err error) bool {
			var status bool
			switch err {
			case Errors.ForceSkipPipelines:
				status = false
			case Errors.ForceSkipMiddlewares:
				status = false
			default:
				status = true
			}

			return status
		}),
	)
	if err != nil {
		catch(err)
	}

}

func (p *BasePipeline) Run(
	context Context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx Context.Context),
) {

	var logCtx interfaces.LoggerInterface
	if ok, logger := log.Get(context); !ok {
		reject(p, Errors.CtxLogFailed)
		l := log.New(log.D{})
		p.Cancel(context, l, Errors.CtxLogFailed)
		return
	} else {
		logCtx = logger
	}

	p.SafeRun(func() error {

		var pError error
		p.run(context, logCtx, func(pipeline interfaces.BasePipelineInterface, err error) {
			pError = err
		}, next)

		return pError

	}, func(err error) {
		reject(p, err)
		p.Cancel(context, logCtx, err)
	})

}

func (p *BasePipeline) Cancel(ctx Context.Context, logger interfaces.LoggerInterface, err error) {
	if utils.IsNill(p.Cn) {
		return
	}
	p.Cn(ctx, logger, err)
}

func (p *BasePipeline) RunMiddlewareStack(
	context Context.Context,
	reject func(middleware interfaces.MiddlewareInterface, err error),
	next func(ctx Context.Context),

) {
	var failed bool
	var forceSkip bool
	var baseException error
	var middlewareContext Context.Context

	middlewareContext = context
	for _, middleware := range p.Middlewares {
		if failed || forceSkip {
			break
		}

		logger := log.New(log.D{"middleware": middleware.GetName(), "pipeline": p.GetName()})

		middleware.Pass(middlewareContext, logger, func(err error) {
			if middleware.IsRequired() {
				baseException = err
				err = Errors.MiddlewareRequired
			}
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
