package net_http

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/elements/forms"
	"net/http"
	"net/url"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
)

func ConstructorProxyMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HttpMiddleware{
		Middleware: &forms.Middleware{
			Required:    required,
			Name:        NAMEHttpBasicAuthMiddleware,
			Description: DESCRIPTIONHttpBasicAuthMiddleware,
		},
		Fn: func(context context.Context, task *task.TyphoonTask, request *http.Request, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {
			transport, ok := ctx.GetContextValue(context,TRANSPORT).(*http.Transport)

			if !ok || !task.Fetcher.IsProxyRequired { reject(Errors.MiddlewareContextFailed); return }

			logrus.Debug("init proxy address ")
			proxyURL, err := url.Parse(task.Fetcher.Proxy)
			if err != nil {
				logrus.Error(err.Error())
				reject(Errors.ProxyUrlWrong)
			}
			if proxyURL.Host != "" && proxyURL.Port() != "" {
				transport.Proxy = http.ProxyURL(proxyURL)
				logrus.Debug("task proxy ", proxyURL.Path)
			} else if proxyURL.Host == "" || proxyURL.Port() == "" {
				reject(Errors.ProxyTaskNotFound)
			}

		},
	}
}
