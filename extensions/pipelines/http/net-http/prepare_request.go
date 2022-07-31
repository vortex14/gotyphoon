package net_http

import (
	"context"
	"net/http"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/interfaces"
)

func CreatePrepareRequestPipeline() *pipelines.TaskPipeline {
	return &pipelines.TaskPipeline{
		BasePipeline: &forms.BasePipeline{
			MetaInfo: &label.MetaInfo{
				Name: "prepare request",
			},
		},
		Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) (error, context.Context) {

			transport, client := GetHttpClientTransport(task)
			request, err := http.NewRequest(task.GetFetcherMethod(), task.GetFetcherUrl(), task.GetRequestBody())

			if err != nil {
				return err, nil
			}

			httpContext := NewClientCtx(context, client)
			httpContext = NewRequestCtx(httpContext, request)
			httpContext = NewTransportCtx(httpContext, transport)

			return nil, httpContext
		},
	}
}
