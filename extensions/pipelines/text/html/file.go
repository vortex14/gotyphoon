package html

import (
	"bytes"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"io/ioutil"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

type FileHtmlPipeline struct {
	*forms.BasePipeline
	*pipelines.TaskPipeline

	Fn func(
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,

		data string,
		doc *goquery.Document,

	) (error, context.Context)

	Cn func(
		err error,
		context context.Context,

		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *FileHtmlPipeline) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
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

	data, err := ioutil.ReadFile(taskInstance.GetFetcherUrl())

	if err != nil {
		reject(t, Errors.FileReadFailed)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(data))
	if err != nil {
		reject(t, err)
		return
	}

	context = NewHtmlCtx(context, doc)

	err, newContext := t.Fn(context, taskInstance, logger, string(data), doc)
	if err != nil {
		reject(t, err)
		return
	}
	next(newContext)
}

func (t *FileHtmlPipeline) Cancel(
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
