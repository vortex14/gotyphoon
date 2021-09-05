package main

import (
	Context "context"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"golang.org/x/net/context"
	"net/http"

	"github.com/vortex14/gotyphoon/data/fake"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func init()  {
	log.InitD()
}

func main() {
	taskTest := fake.CreateDefaultTask()
	taskTest.SetFetcherUrl("http://localhost:12666/fake/product")

	ctxGroup := task.NewTaskCtx(taskTest)


	(&forms.PipelineGroup{
		BaseLabel: interfaces.BaseLabel{
			Name:        "Http strategy",
			Required:    true,
		},
		Stages: []interfaces.BasePipelineInterface{
			&pipelines.TaskPipeline{
				BasePipeline: &forms.BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name: "prepare",
					},
				},
				Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) (error, context.Context) {

					transport, client := net_http.GetHttpClientTransport(task)
					request, err := http.NewRequest(task.GetFetcherMethod(), task.GetFetcherUrl(), nil)

					if err != nil { return err, nil }

					httpContext := net_http.NewClientCtx(context, client)
					httpContext = net_http.NewRequestCtx(httpContext, request)
					httpContext = net_http.NewTransportCtx(httpContext, transport)

					return nil, httpContext
				},
			},
			&net_http.HttpRequestPipeline{
				BasePipeline: &forms.BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name: "http-request",
					},
					Middlewares: []interfaces.MiddlewareInterface{
						//netHttp.ConstructorProxySettingOldMiddleware(true),
						//netHttp.ConstructorProxyRequestSettingsMiddleware(true),
						//netHttp.ConstructorRequestHeaderMiddleware(true),
					},
				},
				Fn: func(context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface, client *http.Client, request *http.Request, transport *http.Transport) (error, Context.Context) {

					err, response, data := net_http.Request(client, request, logger)
					if err != nil { return err, nil }
					context = net_http.NewResponseCtx(context, response)
					context = net_http.NewResponseDataCtx(context, data)


					return nil, context

				},
				Cn: func(err error, context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {
					logger.Error("TEST --- ",err.Error())
				},
			},
			&net_http.HttpResponsePipeline{
				BasePipeline: &forms.BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name: "Response pipeline",
					},
				},
				Fn: func(
					context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface,
					client *http.Client, request *http.Request, transport *http.Transport,
					response *http.Response, data *string) (error, Context.Context){

					logger.Warning("response", *data)

					return nil, context
				},
				Cn: func(err error, context Context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface) {
					logger.Error("pipeline error")
				},
			},
		},
		Consumers:   nil,
	}).Run(ctxGroup)
}
