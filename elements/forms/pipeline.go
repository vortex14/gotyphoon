package forms

import (
	Context "context"
	"errors"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/label"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/utils"
	"golang.org/x/sync/semaphore"

	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

const (
	RetryCount     = "retry_count"
	PanicException = "PANIC"
)

type RetryOptions struct {
	Delay time.Duration

	MaxCount uint

	Required            bool
	OnlyRetryExceptions bool

	RetryExceptions    []error
	CriticalExceptions []error
}

type Options struct {
	Retry         RetryOptions
	MaxConcurrent int64
}

func GetDefaultRetryOptions() *Options {
	return &Options{Retry: RetryOptions{
		MaxCount: 7,
		Delay:    time.Duration(350) * time.Millisecond,
	}}
}

func GetNotRetribleOptions() *Options {
	return &Options{
		Retry: RetryOptions{
			MaxCount: 1,
		},
	}
}

func GetCustomRetryOptions(count uint, delay time.Duration) *Options {
	return &Options{
		Retry: RetryOptions{
			MaxCount: count,
			Delay:    delay,
		},
	}
}

type BasePipeline struct {
	*label.MetaInfo
	awaitabler.Object

	Options        *Options
	NotIgnorePanic bool
	sem            *semaphore.Weighted

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

func (p *BasePipeline) SafeRun(
	context Context.Context,
	logger interfaces.LoggerInterface,
	run func() error, catch func(err error)) {

	if p.Options == nil {
		p.Options = GetDefaultRetryOptions()
	}

	if p.sem == nil && p.Options.MaxConcurrent > 0 {
		p.sem = semaphore.NewWeighted(p.Options.MaxConcurrent)
	}

	if p.sem != nil {
		if !p.sem.TryAcquire(1) {
			logger.Error(Errors.PipelineCrowded)
			catch(Errors.PipelineCrowded)
			return
		}

	}

	defer func() {
		if !p.NotIgnorePanic {
			if r := recover(); r != nil {
				panicE := errors.New(fmt.Sprintf("%s: %s", PanicException, r))
				catch(panicE)
				if p.sem != nil {
					p.sem.Release(1)
				}
			}
		}
	}()
	logger = log.PatchLogI(logger, map[string]interface{}{"max_retry": p.Options.Retry.MaxCount})
	retryCount := 0
	retryMaxCount := p.Options.Retry.MaxCount

	eR := retry.Do(func() error {
		retryCount++
		//logger.Debugf("attempt: %d", retryCount)

		var middlewareErr error

		p.RunMiddlewareStack(context, func(middleware interfaces.MiddlewareInterface, _err error) {
			middlewareErr = _err
			logger.Error("exit from middleware stack . Error: ", middlewareErr.Error())
		}, func(returnedContext Context.Context) {
			context = returnedContext
		})

		if middlewareErr != nil {
			return middlewareErr
		}

		return run()
	},
		retry.Delay(p.Options.Retry.Delay),
		retry.Attempts(retryMaxCount),
		retry.RetryIf(func(_err error) bool {
			var status bool
			switch _err {
			case Errors.ForceSkipPipelines:
				status = false
			case Errors.ForceSkipMiddlewares:
				status = false
			default:
				status = true
			}
			logger.Errorf("RetryIf .... %t, delay: %+v; count: %d", status, p.Options.Retry.Delay, retryCount)
			return status
		}),
	)

	if eR != nil {
		catch(eR)
	}

	if p.sem != nil {
		p.sem.Release(1)
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
	var pError error

	p.SafeRun(context, logCtx, func() error {

		if utils.IsNill(p.Fn) {
			return Errors.LambdaRequired
		}

		err, newContext := p.Fn(context, logCtx)
		if utils.NotNill(err) {
			return err
		}

		next(newContext)

		return nil

	}, func(err error) {
		pError = err
		reject(p, pError)
		p.Cancel(context, logCtx, pError)
	})

}

func (p *BasePipeline) Cancel(ctx Context.Context, logger interfaces.LoggerInterface, err error) {
	if p.Cn == nil {
		return
	}
	logger.Error(p.Cn, utils.IsNill(p.Cn))
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
		logger.Debug("Run")
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

func InsertPipeline(
	a []interfaces.BasePipelineInterface,
	index int, value interfaces.BasePipelineInterface) []interfaces.BasePipelineInterface {

	if len(a) == index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}
