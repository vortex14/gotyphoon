package rod

import (
	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/task"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/extensions/pipelines/text/html"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type HttpRodResponsePipeline struct {
	*forms.BasePipeline
	*pipelines.TaskPipeline

	Fn func(
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,

		browser *rod.Browser,
		page *rod.Page,
		body *string,
		doc *goquery.Document,

	) (error, context.Context)

	Cn func(
		err error,
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *HttpRodResponsePipeline) UnpackResponseCtx(
	ctx context.Context,
) (bool, interfaces.TaskInterface, interfaces.LoggerInterface, *rod.Browser, *rod.Page, *string, *goquery.Document) {
	okT, taskInstance := task.Get(ctx)
	okL, logger := log.Get(ctx)

	okB, browser := GetBrowserCtx(ctx)
	okP, page := GetPageCtx(ctx)
	okR, body := GetPageResponse(ctx)
	okD, doc := html.GetHtmlDoc(ctx)

	if !okT || !okL || !okB || !okP || !okR || !okD {
		return false, nil, nil, nil, nil, nil, nil
	}

	return okL && okT && okB, taskInstance, logger, browser, page, body, doc
}

func (t *HttpRodResponsePipeline) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
) {

	if t.Fn == nil {
		reject(t, Errors.TaskPipelineRequiredHandler)
		return
	}

	t.SafeRun(func() error {
		ok, taskInstance, logger, browser, page, body, doc := t.UnpackResponseCtx(context)

		if !ok {
			return fmt.Errorf("%s. taskInstance: %v, logger: %v, browser: %v, page: %v, body: %v",
				Errors.PipelineContexFailed, taskInstance, logger, browser, page, body)
		}

		err, newContext := t.Fn(context, taskInstance, logger, browser, page, body, doc)
		if err != nil {
			return err
		}
		next(newContext)
		return err

	}, func(err error) {
		reject(t, err)
		_, logCtx := log.Get(context)
		t.Cancel(context, logCtx, err)
	})

}

func (t *HttpRodResponsePipeline) Cancel(
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
