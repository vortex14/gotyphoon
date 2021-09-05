package net_http

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/interfaces"
	"net/http"
)

const (
	NAME = "Request Header Middleware"
	DESCRIPTION = "Setting request header for Typhoon task"
)

func ConstructorRequestHeaderMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HttpMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required:    required,
				Name:        NAME,
				Description: DESCRIPTION,
			},
		},
		Fn: func(context context.Context, task *task.TyphoonTask, request *http.Request,
			logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context),
		) {
			for key, element := range task.Fetcher.Headers {
				request.Header.Add(
					key,
					element,
				)
			}
		},
	}
}
