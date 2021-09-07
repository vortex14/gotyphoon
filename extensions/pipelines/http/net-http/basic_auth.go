package net_http

import (
	"context"
	"net/http"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	BASICAuthUserTaskKey = "login"
	BASICAuthPasswordTaskKey = "password"

	NAMEHttpBasicAuthMiddleware = "Http Basic auth middleware"
	DESCRIPTIONHttpBasicAuthMiddleware = "Setting basic auth credentials for http request from Typhoon task"
)

func ConstructorBasicAuthMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HttpMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required:    required,
				Name:        NAMEHttpBasicAuthMiddleware,
				Description: DESCRIPTIONHttpBasicAuthMiddleware,
			},
		},

		Fn: func(
			context context.Context, task *task.TyphoonTask, request *http.Request,
			logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context),

			) {

			login, okL := task.Fetcher.Auth[BASICAuthUserTaskKey]
			passwd, okP := task.Fetcher.Auth[BASICAuthPasswordTaskKey]

			if okL && okP { request.SetBasicAuth(login, passwd) }
		},
	}
}