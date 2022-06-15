package gin

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/vortex14/gotyphoon/elements/forms"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type RequestPipeline struct {
	*forms.BasePipeline

	Fn func(
		context context.Context,
		ginCtx *gin.Context,
		logger interfaces.LoggerInterface,
	)  (error, context.Context)

	Cn func(
		err error,
		context context.Context,
		ginCtx *gin.Context,
		logger interfaces.LoggerInterface,
	)
}

func (t *RequestPipeline) UnpackCtx(ctx context.Context) (bool, *gin.Context, interfaces.LoggerInterface)  {
	okG, ginCtx := GetRequestCtx(ctx)
	okL, logger := log.Get(ctx)
	return okL && okG, ginCtx, logger
}

func (t *RequestPipeline) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
) {

	if t.Fn == nil { reject(t, Errors.TaskPipelineRequiredHandler); return }

	ok,ginCtx, logger := t.UnpackCtx(context)
	if !ok { reject(t, Errors.PipelineContexFailed); return }

	err, newContext := t.Fn(context, ginCtx, logger)
	if err != nil { reject(t, err); return }
	next(newContext)
}

func (t *RequestPipeline) Cancel(
	context context.Context,
	logger interfaces.LoggerInterface,
	err error,
) {

	if t.Cn == nil { return }

	ok,ginCtx, logger := t.UnpackCtx(context)
	if !ok { return }

	t.Cn(err, context, ginCtx, logger)

}
