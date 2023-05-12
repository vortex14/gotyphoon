package html

import (
	"bytes"
	Context "context"
	"github.com/vortex14/gotyphoon/log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/pipelines"

	Errors "github.com/vortex14/gotyphoon/errors"
	netHttp "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/interfaces"
)

type ResponseHtmlPipeline struct {
	*forms.BasePipeline
	*pipelines.TaskPipeline

	netHttp.HttpResponsePipeline

	Fn func(
		context Context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,

		request *http.Request,
		response *http.Response,
		data *string,
		doc *goquery.Document,

	) (error, Context.Context)

	Cn func(
		err error,
		context Context.Context,

		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *ResponseHtmlPipeline) Run(
	context Context.Context,
	reject func(context Context.Context, pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx Context.Context),
) {

	if t.Fn == nil {
		reject(context, t, Errors.TaskPipelineRequiredHandler)
		return
	}

	_, logger := log.Get(context)

	t.SafeRun(context, logger, func(patchedContext Context.Context) error {

		ok, taskInstance, logger, _, request, _, response, data := t.UnpackResponse(patchedContext)

		if !ok {
			return Errors.PipelineContexFailed
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(*data)))
		if err != nil {
			return err
		}

		patchedContext = NewHtmlCtx(patchedContext, doc)

		err, newContext := t.Fn(patchedContext, taskInstance, logger, request, response, data, doc)
		if err != nil {
			return err
		}
		next(newContext)
		return nil

	}, func(context Context.Context, err error) {
		reject(context, t, err)
		t.Cancel(context, logger, err)
	})

}

func (t *ResponseHtmlPipeline) Cancel(
	context Context.Context,
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
