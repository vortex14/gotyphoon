package rod

import (
	"context"
	"fmt"

	"github.com/go-rod/rod"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/task"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type HttpRodRequestPipeline struct {
	*forms.BasePipeline
	*pipelines.TaskPipeline

	Fn func(
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,

		browser *rod.Browser,

	) (error, context.Context)

	Cn func(
		err error,
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *HttpRodRequestPipeline) UnpackRequestCtx(
	ctx context.Context,
) (bool, interfaces.TaskInterface, interfaces.LoggerInterface, *rod.Browser) {
	okT, taskInstance := task.Get(ctx)
	okL, logger := log.Get(ctx)

	okB, browser := GetBrowserCtx(ctx)

	if !okT || !okL || !okB {
		return false, nil, nil, nil
	}

	return okL && okT && okB, taskInstance, logger, browser
}

func (t *HttpRodRequestPipeline) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
) {

	if t.Fn == nil {
		reject(t, Errors.TaskPipelineRequiredHandler)
		return
	}

	t.SafeRun(func() error {
		ok, taskInstance, logger, browser := t.UnpackRequestCtx(context)

		if !ok {
			return fmt.Errorf("%s. taskInstance: %v, logger: %v, browser: %v", Errors.PipelineContexFailed, taskInstance, logger, browser)
		}

		err, newContext := t.Fn(context, taskInstance, logger, browser)
		if err != nil {
			return err
		}
		next(newContext)
		return err

	}, func(err error) {

		// without this will be leaked after panic.
		if e := rod.Try(func() {
			_, b := GetBrowserCtx(context)
			b.MustClose()
		}); e != nil {
			reject(t, e)
			_, logCtx := log.Get(context)
			t.Cancel(context, logCtx, e)

			return
		}

		reject(t, err)
		_, logCtx := log.Get(context)
		t.Cancel(context, logCtx, err)
	})

}

func (t *HttpRodRequestPipeline) Cancel(
	context context.Context,
	logger interfaces.LoggerInterface,
	err error,
) {

	if t.Cn == nil {
		return
	}

	ok, taskInstance, logger := t.UnpackCtx(context)
	if !ok {
		return
	}

	t.Cn(err, context, taskInstance, logger)

}
