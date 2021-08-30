package pipelines

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/forms"

	"github.com/vortex14/gotyphoon/ctx"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
)

type TaskPipeline struct {
	Fn func(context context.Context, task *task.TyphoonTask, logger interfaces.LoggerInterface)  (error, context.Context)
	*forms.BasePipeline
}

func (t *TaskPipeline) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
	) {

	if t.Fn == nil { reject(t, Errors.TaskPipelineRequiredHandler); return }

	taskInstance, okT := ctx.GetContextValue(context, interfaces.TASK).(*task.TyphoonTask)
	logger,okL := ctx.GetContextValue(context, interfaces.LOGGER).(interfaces.LoggerInterface)
	if !okL || !okT { reject(t, Errors.PipelineContexFailed); return }

	err, newContext := t.Fn(context, taskInstance, logger)
	if err != nil { reject(t, err); return }
	next(newContext)
}