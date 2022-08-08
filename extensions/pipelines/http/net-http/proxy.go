package net_http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fatih/color"

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
)

func ConstructorProxySettingMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HttpMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required:    required,
				Name:        NAMEProxyMiddleware,
				Description: DescriptionProxyMiddleware,
			},
		},
		Fn: func(context context.Context, task *task.TyphoonTask, request *http.Request, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {

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

			urlSupported := fmt.Sprintf("%s/proxy", proxyServer)
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
				color.Red("JsonLoad has Error: %s", err.Error())
				reject(err)
				return
			}
			if !proxyResponse.Success {
				reject(Errors.ProxyBusy)
				return
			}
			task.SetUserAgent(proxyResponse.Agent)
			task.SetProxyAddress(proxyResponse.Proxy)

		},
	}
}
