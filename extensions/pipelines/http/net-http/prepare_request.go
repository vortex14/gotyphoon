package net_http

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/interfaces"
	"net/http"
)

func CreatePrepareRequestPipeline() *pipelines.TaskPipeline {
	return &pipelines.TaskPipeline{
		BasePipeline: &forms.BasePipeline{
			Name: "prepare",
		},
		Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) (error, context.Context) {

			transport, client := GetHttpClientTransport(task)
			request, err := http.NewRequest(task.GetFetcherMethod(), task.GetFetcherUrl(), nil)

			if err != nil { return err, nil }

			httpContext := NewClientCtx(context, client)
			httpContext = NewRequestCtx(httpContext, request)
			httpContext = NewTransportCtx(httpContext, transport)

			return nil, httpContext
		},
	}
}
