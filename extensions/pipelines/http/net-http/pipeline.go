package net_http

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/log"
	"net/http"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

type HttpRequestPipeline struct {
	*forms.BasePipeline
	*pipelines.TaskPipeline
	Fn func(
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,

		client *http.Client,
		request *http.Request,
		transport *http.Transport,

	)  (error, context.Context)
	Cn func(err error, context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface)
}

func (t *HttpRequestPipeline) UnpackRequestCtx(ctx context.Context) (bool, interfaces.TaskInterface, interfaces.LoggerInterface, *http.Client, *http.Request, *http.Transport) {
	okT, taskInstance := task.Get(ctx)
	okL, logger := log.Get(ctx)

	okC, client := GetClientCtx(ctx)
	okR, request := GetRequestCtx(ctx)
	okTr, transport := GetTransportCtx(ctx)

	if (!okC || !okR) || (!okT || !okL || !okTr) { return false, nil, nil, nil, nil, nil }

	return okL && okT && okC && okR && okTr, taskInstance, logger, client, request, transport
}

func (t *HttpRequestPipeline) Run(
	context context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx context.Context),
) {

	if t.Fn == nil { reject(t, Errors.TaskPipelineRequiredHandler); return }

	ok, taskInstance, logger, client, request, transport := t.UnpackRequestCtx(context)

	if !ok { reject(t, Errors.PipelineContexFailed); return }

	err, newContext := t.Fn(context, taskInstance, logger, client, request, transport)
	if err != nil { reject(t, err); return }
	next(newContext)
}

func (t *HttpRequestPipeline) Cancel(
	context context.Context,
	logger interfaces.LoggerInterface,
	err error,
) {

	if t.Fn == nil { return }

	ok,taskInstance, logger := t.UnpackCtx(context)
	if !ok { return }

	t.Cn(err, context, taskInstance, logger)

}