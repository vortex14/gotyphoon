package middlewares

import (
	"context"

	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
)

type MiddlewareTaskCallback func(
	context context.Context,
	task *task.TyphoonTask,
	logger interfaces.LoggerInterface,
	reject func(err error),
	next func(ctx context.Context),
)

type TaskMiddleware struct {
	*forms.Middleware
	Fn MiddlewareTaskCallback
}

func (m *TaskMiddleware) Pass(
	context context.Context,
	logger interfaces.LoggerInterface,
	reject func(err error),
	next func(ctx context.Context),

	) {

	taskInstance, _ := ctx.Get(context, interfaces.TASK).(*task.TyphoonTask)
	m.Fn(context, taskInstance, logger, reject, next)
}


