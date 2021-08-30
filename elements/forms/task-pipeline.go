package forms

import (
	"context"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
)

type TaskPipeline struct {
	TaskHandler   func(context context.Context, task *task.TyphoonTask, logger interfaces.LoggerInterface)  (error, context.Context)
	*BasePipeline
}

func (t *TaskPipeline) Run(ctx context.Context) (error, context.Context) {
	if t.TaskHandler == nil {
		return Errors.TaskPipelineRequiredHandler, nil
	}
	taskInstance, _ := ctx.Value(interfaces.ContextKey(interfaces.TASK)).(*task.TyphoonTask)
	logger,_ := ctx.Value(interfaces.ContextKey(interfaces.LOGGER)).(interfaces.LoggerInterface)
	return t.TaskHandler(ctx, taskInstance, logger)
}