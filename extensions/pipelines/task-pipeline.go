package pipelines

import (
	Context "context"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/task"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type TaskPipeline struct {
	*forms.BasePipeline

	Fn func(
		context Context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	) (error, Context.Context)

	Cn func(
		err error,
		context Context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *TaskPipeline) UnpackCtx(
	ctx Context.Context,
) (bool, interfaces.TaskInterface, interfaces.LoggerInterface) {

	okT, taskInstance := task.Get(ctx)
	okL, logger := log.Get(ctx)
	return okL && okT, taskInstance, logger
}

func (t *TaskPipeline) Run(
	context Context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx Context.Context),
) {

	if t.Fn == nil {
		reject(t, Errors.TaskPipelineRequiredHandler)
		return
	}

	ok, taskInstance, logger := t.UnpackCtx(context)
	if !ok {
		reject(t, Errors.PipelineContexFailed)
		return
	}

	t.SafeRun(context, logger, func(patchedCtx Context.Context) error {
		err, newContext := t.Fn(context, taskInstance, logger)
		if err != nil {
			return err
		}

		next(newContext)

		return nil

	}, func(err error) {
		reject(t, err)
		t.Cancel(context, logger, err)
	})

}

func (t *TaskPipeline) Cancel(
	context Context.Context,
	logger interfaces.LoggerInterface,
	err error,
) {

	if t.Cn == nil {
		return
	}

	ok, taskInstance, _logger := t.UnpackCtx(context)
	if !ok {
		logger.Error(Errors.PipelineContexFailed)
		return
	}

	t.Cn(err, context, taskInstance, _logger)

}
