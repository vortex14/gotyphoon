package interfaces

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/task"
	"sync"
	"time"
)



type ResponseInterface interface {

}


type HandlerInterface interface {

}


type BasePipelineInterface interface {
	Run() (error, interface{})

	NextStage()
	Finish()
	Retry()
	Await()
}

type ProcessorPipelineInterface interface {
	BasePipelineInterface
	Crawl()
	Switch()
}

type PipelineMiddleware interface {

}

type CallbackPipelineInterface interface {
	Call(ctx context.Context, data interface{})
}

type BasePipelineLabel struct {
	Name        string
	Description string
}

type ConsumerInterface interface {

}

type LambdaInterface interface {

}

type BasePipeline struct {
	Task          *task.TyphoonTask
    Project       Project

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


	Handler     HandlerInterface
	Response    ResponseInterface

	LOG         *logrus.Entry
	Metrics     MetricsInterface

	Middlewares []PipelineMiddleware

	Callbacks   []CallbackPipelineInterface
	Consumers   map[string]ConsumerInterface

	LambdaHandler func(ctx context.Context, data interface{}) error
	CancelHandler func(ctx context.Context, data interface{}, err error)

	*BasePipelineLabel
}

func (p *BasePipeline) NextStage()  {

}

func (p *BasePipeline) Await()  {
	p.promise.Wait()
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


type PipelineGroup struct {
	*BasePipelineLabel

	errorCount    int64
	duration      time.Time
	timeLife      time.Time

	LambdaMap     map[string]LambdaInterface
	PyLambdaMap   map[string]LambdaInterface

	Stages      []BasePipelineInterface
	Consumers   map[string]ConsumerInterface

}