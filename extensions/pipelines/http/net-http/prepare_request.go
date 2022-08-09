package net_http

import (
	"context"
	"net/http"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/middlewares"
	"github.com/vortex14/gotyphoon/interfaces"
)

func ConstructorPrepareRequestMiddleware(required bool) interfaces.MiddlewareInterface {
	return &middlewares.TaskMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required: required,
				Name:     "prepare request",
			},
		},
		Fn: func(context context.Context, task *task.TyphoonTask,
			logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {

			transport, client := GetHttpClientTransport(task)
			request, err := http.NewRequest(task.GetFetcherMethod(), task.GetFetcherUrl(), task.GetRequestBody())

			if err != nil {
				reject(err)
				return
			}

			httpContext := NewClientCtx(context, client)
			httpContext = NewRequestCtx(httpContext, request)
			httpContext = NewTransportCtx(httpContext, transport)

			next(httpContext)

		},
	}
}
