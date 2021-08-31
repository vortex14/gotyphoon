package pipelines

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/log"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
)

type TaskPipeline struct {
	Fn func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface)  (error, context.Context)
	Cn func(err error, context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface)
	*forms.BasePipeline
}

func (t *TaskPipeline) UnpackCtx(ctx context.Context) (bool, interfaces.TaskInterface, interfaces.LoggerInterface)  {
	okT, taskInstance := task.Get(ctx)
	okL, logger := log.Get(ctx)
	return okL && okT, taskInstance, logger
}

func (t *TaskPipeline) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
	) {

	if t.Fn == nil { reject(t, Errors.TaskPipelineRequiredHandler); return }

	ok,taskInstance, logger := t.UnpackCtx(context)
	if !ok { reject(t, Errors.PipelineContexFailed); return }

	err, newContext := t.Fn(context, taskInstance, logger)
	if err != nil { reject(t, err); return }
	next(newContext)
}

func (t *TaskPipeline) Cancel(
	context context.Context,
	logger interfaces.LoggerInterface,
	err error,
) {

	if t.Cn == nil { return }

	ok,taskInstance, logger := t.UnpackCtx(context)
	if !ok { return }

	t.Cn(err, context, taskInstance, logger)

}