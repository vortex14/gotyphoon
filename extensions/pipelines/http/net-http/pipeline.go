package net_http

import (
	Context "context"
	"golang.org/x/net/context"
	"net/http"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/task"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type HttpRequestPipeline struct {
	*forms.BasePipeline
	*pipelines.TaskPipeline

	Fn func(
		context Context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,

		client *http.Client,
		request *http.Request,
		transport *http.Transport,

	) (error, Context.Context)

	Cn func(
		err error,
		context context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *HttpRequestPipeline) UnpackRequestCtx(
	ctx Context.Context,
) (bool, interfaces.TaskInterface, interfaces.LoggerInterface, *http.Client, *http.Request, *http.Transport) {
	okT, taskInstance := task.Get(ctx)
	okL, logger := log.Get(ctx)

	okC, client := GetClientCtx(ctx)
	okR, request := GetRequestCtx(ctx)
	okTr, transport := GetTransportCtx(ctx)

	if (!okC || !okR) || (!okT || !okL || !okTr) {
		logger.Errorf("client: %t , request: %t, transport: %t", okC, okR, okTr)
		return false, nil, nil, nil, nil, nil
	}

	return okL && okT && okC && okR && okTr, taskInstance, logger, client, request, transport
}

func (t *HttpRequestPipeline) Run(
	context Context.Context,
	reject func(context Context.Context, pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx Context.Context),
) {

	if t.Fn == nil {
		reject(context, t, Errors.TaskPipelineRequiredHandler)
		return
	}

	_, logger := log.Get(context)

	t.SafeRun(context, logger, func(patchedCtx Context.Context) error {

		ok, taskInstance, logger, client, request, transport := t.UnpackRequestCtx(context)

		if !ok {
			return Errors.PipelineContexFailed
		}

		err, newContext := t.Fn(context, taskInstance, logger, client, request, transport)
		if err != nil {
			return err
		}
		next(newContext)
		return err

	}, func(context Context.Context, err error) {
		reject(context, t, err)
		_, logCtx := log.Get(context)
		t.Cancel(context, logCtx, err)
	})

}

func (t *HttpRequestPipeline) Cancel(
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
