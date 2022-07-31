package html

import (
	"bytes"
	"context"
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
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,

		request *http.Request,
		response *http.Response,
		data *string,
		doc *goquery.Document,

	) (error, context.Context)

	Cn func(
		err error,
		context context.Context,

		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *ResponseHtmlPipeline) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
) {

	if t.Fn == nil {
		reject(t, Errors.TaskPipelineRequiredHandler)
		return
	}

	ok, taskInstance, logger, _, request, _, response, data := t.UnpackResponse(context)

	if !ok {
		reject(t, Errors.PipelineContexFailed)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(*data)))
	if err != nil {
		reject(t, err)
		return
	}

	context = NewHtmlCtx(context, doc)

	err, newContext := t.Fn(context, taskInstance, logger, request, response, data, doc)
	if err != nil {
		reject(t, err)
		return
	}
	next(newContext)
}

func (t *ResponseHtmlPipeline) Cancel(
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
