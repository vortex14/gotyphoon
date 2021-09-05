package net_http

import (
	"context"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

const (
	InstallingProxyMiddleware = "Proxy middleware"
	DescriptionDProxyMiddleware = "Proxy Middleware"
)

func ConstructorProxyRequestSettingsMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HttpMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required:    required,
				Name:        InstallingProxyMiddleware,
				Description: DescriptionDProxyMiddleware,
			},
		},
		Fn: func(context context.Context,
			task *task.TyphoonTask,
			request *http.Request,
			logger interfaces.LoggerInterface,
			reject func(err error),
			next func(ctx context.Context)) {

			ok, transport := GetTransportCtx(context)
			if !ok { reject(Errors.MiddlewareContextFailed); return }
			if !task.IsProxyRequired() { reject(Errors.ProxyTaskRequired); return}
			logrus.Debug("init proxy address ...", task.GetProxyAddress())
			proxyURL, err := url.Parse(task.GetProxyAddress())
			if err != nil || proxyURL == nil {
				logrus.Error(err.Error())
				reject(Errors.ProxyUrlWrong)
				return
			}
			if proxyURL.Host != "" && proxyURL.Port() != "" {
				transport.Proxy = http.ProxyURL(proxyURL)
			} else if proxyURL.Host == "" || proxyURL.Port() == "" {
				reject(Errors.ProxyTaskNotFound)
			}
		},
	}
}
