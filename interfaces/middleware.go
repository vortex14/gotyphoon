package interfaces

import (
	"context"
	"github.com/vortex14/gotyphoon/task"
)

type ConstructorMiddleware func (required bool) MiddlewareInterface

type MiddlewareCallback func(
		context context.Context,
		logger LoggerInterface,
		reject func(err error),
		next func(ctx context.Context),
	)

type MiddlewareTaskCallback func(
		context context.Context,
		task *task.TyphoonTask,
		logger LoggerInterface,
		reject func(err error),
		next func(ctx context.Context),
	)

type MiddlewareInterface interface {
	IsRequired() bool
	Pass(context context.Context,
		logger LoggerInterface,
		reject func(err error),
		next func(context context.Context),
	)

	MetaDataInterface
}

type Middleware struct {
	Name        string
	Required    bool
	Description string
	Callback    MiddlewareCallback
	PyCallback  MiddlewareCallback
}

func (m *Middleware) GetName() string {
	return m.Name
}

func (m *Middleware) GetDescription() string {
	return m.Description
}

func (m *Middleware) IsRequired() bool {
	return m.Required
}

func (m *Middleware) Pass(context context.Context, logger LoggerInterface, reject func(err error), next func(ctx context.Context)) {
	m.Callback(context, logger, reject, next)
}

type TaskMiddleware struct {
	*Middleware
	Callback MiddlewareTaskCallback
}

func (m *TaskMiddleware) Pass(ctx context.Context, logger LoggerInterface, reject func(err error), next func(ctx context.Context)) {
	taskInstance, _ := ctx.Value(TASK).(*task.TyphoonTask)
	m.Callback(ctx, taskInstance, logger, reject, next)
}