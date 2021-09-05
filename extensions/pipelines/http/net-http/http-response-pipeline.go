package net_http

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"net/http"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

type HttpResponsePipeline struct {
	*forms.BasePipeline
	*pipelines.TaskPipeline
	*HttpRequestPipeline

	Fn func(
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,

		client *http.Client,
		request *http.Request,
		transport *http.Transport,

		response *http.Response,
		data *string,

	)  (error, context.Context)

	Cn func(
		err error,
		context context.Context,

		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *HttpResponsePipeline) UnpackResponse(ctx context.Context) (
	bool,
	interfaces.TaskInterface,
	interfaces.LoggerInterface,

	*http.Client,
	*http.Request,
	*http.Transport,

	*http.Response,
	*string,
	) {

	ok,taskInstance, logger, client, request, transport := t.UnpackRequestCtx(ctx)

	okR, response, data := GetResponseCtx(ctx)
	if !ok || !okR { return false, nil, nil, nil, nil, nil, nil, nil }
	return ok, taskInstance, logger, client, request, transport, response, data
}

func (t *HttpResponsePipeline) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
) {

	if t.Fn == nil { reject(t, Errors.TaskPipelineRequiredHandler); return }

	ok, taskInstance, logger, client, request, transport, response, data := t.UnpackResponse(context)

	if !ok { reject(t, Errors.PipelineContexFailed); return }

	err, newContext := t.Fn(context, taskInstance, logger, client, request, transport, response, data)
	if err != nil { reject(t, err); return }
	next(newContext)
}

func (t *HttpResponsePipeline) Cancel(
	context context.Context,
	logger interfaces.LoggerInterface,
	err error,
) {

	if t.Cn == nil { return }
	ok, taskInstance, logger := t.UnpackCtx(context)
	if !ok { return }
	t.Cn(err, context, taskInstance, logger)
}