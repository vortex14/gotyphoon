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
	NAMEProxyOldMiddleware = "Installation proxy middleware"
	DescriptionProxyOldMiddleware = "Setting proxy Middleware"
)

func ConstructorProxySettingOldMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HttpMiddleware{
		Middleware: &forms.Middleware{
			MetaInfo: &label.MetaInfo{
				Required:    required,
				Name:        NAMEProxyOldMiddleware,
				Description: DescriptionProxyOldMiddleware,
			},
		},
		Fn: func(context context.Context, task *task.TyphoonTask, request *http.Request, logger interfaces.LoggerInterface, reject func(err error), next func(ctx context.Context)) {


			client := GetHttpClient(task)

			urlOldSupported := fmt.Sprintf("%s/url=%s", task.GetProxyServerUrl(), task.GetFetcherUrl())
			logger.Info(urlOldSupported)

			request, err := http.NewRequest(http.MethodGet, urlOldSupported, nil)

			if err != nil { reject(err); return }

			err, body, _ := GetBody(client, request)
			if err != nil || body == nil {reject(Errors.ProxyServerError); return}
			logger.Error(string(*body))
			var proxyResponse models.Proxy
			err = utils.JsonLoad(&proxyResponse, *body)
			if err != nil {
				color.Red("JsonLoad has Error: %s", err.Error())
				reject(err)
				return
			}
			proxyFormat := FormattingProxy(proxyResponse.Proxy)
			task.SetUserAgent(proxyResponse.Agent)
			task.SetProxyAddress(proxyFormat)

		},
	}
}
