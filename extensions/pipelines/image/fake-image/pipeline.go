package fake_image

import (
	"context"
	"github.com/fogleman/gg"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/task"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type ImagePipeline struct {
	*forms.BasePipeline

	Fn func(
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
		imgCtx *gg.Context,
	) (error, context.Context)

	Cn func(
		err error,
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *ImagePipeline) UnpackCtx(
	ctx context.Context,
) (bool, interfaces.TaskInterface, interfaces.LoggerInterface, *gg.Context) {

	okT, taskInstance := task.Get(ctx)
	okL, logger := log.Get(ctx)
	okD, data := GetImgCtx(ctx)
	return okL && okT && okD, taskInstance, logger, data
}

func (t *ImagePipeline) Run(
	context context.Context,
	reject func(context context.Context, pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
) {

	if t.Fn == nil {
		reject(context, t, Errors.TaskPipelineRequiredHandler)
		return
	}

	ok, taskInstance, logger, imgCtx := t.UnpackCtx(context)

	if !ok {
		reject(context, t, Errors.PipelineContexFailed)
		return
	}

	err, newContext := t.Fn(context, taskInstance, logger, imgCtx)
	if err != nil {
		reject(context, t, err)
		return
	}
	next(newContext)
}

func (t *ImagePipeline) Cancel(
	context context.Context,
	logger interfaces.LoggerInterface,
	err error,
) {

	if t.Cn == nil {
		return
	}

	ok, taskInstance, logger, _ := t.UnpackCtx(context)
	if !ok {
		return
	}

	t.Cn(err, context, taskInstance, logger)

}
