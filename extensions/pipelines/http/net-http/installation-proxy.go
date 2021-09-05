package net_http

import (
	"context"
	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"net/http"

	"github.com/vortex14/gotyphoon/elements/forms"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/extensions/models"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)

const (
	NAMEProxyMiddleware = "Installation proxy middleware"
	DescriptionProxyMiddleware = "Setting proxy Middleware"
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

			client := GetHttpClient(task)

			request, err := http.NewRequest(task.GetFetcherMethod(), task.GetFetcherUrl(), nil)

			if err != nil { reject(err); return }

			err, body, _ := GetBody(client, request)
			if err != nil || body == nil {reject(Errors.ProxyServerError); return}

			var proxyResponse models.Proxy
			err = utils.JsonLoad(&proxyResponse, *body)
			if err != nil {
				color.Red("%s", err.Error())
				reject(err)
				return 
			}

			task.SetUserAgent(proxyResponse.Agent)
			task.SetProxyAddress(proxyResponse.Proxy)

		},
	}
}
