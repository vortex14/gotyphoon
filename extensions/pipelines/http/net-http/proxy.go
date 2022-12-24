package net_http

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vortex14/gotyphoon/extensions/middlewares"
	"github.com/vortex14/gotyphoon/log"
	"net/http"
	"net/url"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/extensions/models"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)

const (
	NAMEProxyMiddleware        = "Installation proxy pipeline"
	DescriptionProxyMiddleware = "Setting proxy"

	InstallingProxyMiddleware   = "Proxy middleware"
	DescriptionDProxyMiddleware = "Proxy Middleware"
)

func ConstructorProxySettingMiddleware(required bool) interfaces.MiddlewareInterface {
	return &middlewares.TaskMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required: required,
				Name:     "get proxy for request",
			},
		},
		Fn: func(context context.Context, task *task.TyphoonTask,
			logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {

			if len(task.GetProxyAddress()) > 0 {
				next(context)
				return
			}

			client := GetHttpClient(task)

			proxyServer := task.GetProxyServerUrl()

			if len(proxyServer) == 0 {
				reject(Errors.ProxyServerNotFound)
				return
			}

			urlSupported := fmt.Sprintf("%s/proxy?url=%s&encode=base64", proxyServer, task.GetBase64FetcherURL())
			logger.Info("get proxy from :", urlSupported)

			request, err := http.NewRequest(http.MethodGet, urlSupported, nil)

			if err != nil {
				reject(err)
				return
			}

			err, body, _ := GetBody(client, request)
			if err != nil || body == nil {
				reject(Errors.ProxyServerError)
				return
			}
			//logger.Error(string(*body))
			var proxyResponse models.Proxy
			err = utils.JsonLoad(&proxyResponse, *body)
			if err != nil {
				logger.Debug(fmt.Sprintf("JsonLoad has Error: %s", err.Error()))
				reject(err)
				return
			}
			if !proxyResponse.Success {
				reject(Errors.ProxyBusy)
				return
			}
			task.SetUserAgent(proxyResponse.Agent)
			task.SetProxyAddress(proxyResponse.Proxy)

			next(context)

		},
	}
}

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
			if !ok {
				reject(Errors.MiddlewareContextFailed)
				return
			}
			if !task.IsProxyRequired() && len(task.GetProxyAddress()) == 0 {
				reject(Errors.ProxyTaskRequired)
				return
			}
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

			context = log.PatchCtx(context, map[string]interface{}{"proxy": task.GetProxyAddress()})
			next(context)

		},
	}
}
