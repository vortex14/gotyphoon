package net_http

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/extensions/middlewares"
	"github.com/vortex14/gotyphoon/interfaces"
)

func ConstructorMockResponseMiddleware(required bool) interfaces.MiddlewareInterface {
	return &middlewares.TaskMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required: required,
				Name:     "mock response",
			},
		},
		Fn: func(context context.Context, task *task.TyphoonTask,
			logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
			bodyStr := "{\"s\":\"test_json\"}"
			r := io.NopCloser(strings.NewReader(bodyStr)) // r type is io.ReadCloser
			response := http.Response{
				StatusCode: 200,
				Body:       r,
			}
			httpContext := NewResponseCtx(context, &response)

			httpContext = NewResponseDataCtx(httpContext, &bodyStr)

			next(httpContext)

		},
	}
}
