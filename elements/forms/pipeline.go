package forms

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"

)

type BasePipeline struct {
	Task          *task.TyphoonTask
	Project       interfaces.Project

	Context       context.Context
	stageIndex    int32
	promise       sync.WaitGroup

	inputCount    int64
	inputByte     int64
	InputMap      map[string]interface{}

	outputCount   int64
	outputByte    int64
	OutputMap     map[string]interface{}


	errorCount    int64
	duration      time.Time


	Handler     interfaces.HandlerInterface
	Response    interfaces.ResponseInterface

	LOG         *logrus.Entry
	Metrics     interfaces.MetricsInterface

	Middlewares []interfaces.MiddlewareInterface

	Callbacks   []interfaces.CallbackPipelineInterface
	Consumers   map[string]interfaces.ConsumerInterface

	LambdaHandler func(ctx context.Context, data interface{}) error
	CancelHandler func(ctx context.Context, data interface{}, err error)

	*interfaces.BaseLabel
}

func (p *BasePipeline) NextStage()  {

}

func (p *BasePipeline) Await()  {
	p.promise.Wait()
}

func (p *BasePipeline) IsRequired() bool {
	return p.Required
}

func (p *BasePipeline) Run(ctx context.Context, data interface{}) error {
	return p.LambdaHandler(ctx, data)
}

func (p *BasePipeline) Cancel(data interface{}, err error) {
	p.CancelHandler(p.Context, data, err)
}

func (p *BasePipeline) Wrap(
	lambda func(ctx context.Context, data interface{},
) error) {
	p.LambdaHandler = lambda
}

func (p *BasePipeline) RunMiddlewareStack(
	context context.Context,
	reject func(middleware interfaces.MiddlewareInterface,err error),

	) {
	var failed bool
	for _, middleware := range p.Middlewares {
		if failed {break}

		middlewareLogger := logrus.WithFields(logrus.Fields{
			"middleware": middleware.GetName(),
			"pipeline": p.GetName(),
		})

		middleware.Pass(context, middlewareLogger, func(err error) {
			if middleware.IsRequired() {
				reject(middleware, err)
				failed = true
			} else {
				middlewareLogger.Warning(err.Error())
			}
		})
	}
}
