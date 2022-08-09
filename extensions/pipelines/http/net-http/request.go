package net_http

import (
	"context"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/vortex14/gotyphoon/extensions/models"
	"github.com/vortex14/gotyphoon/utils"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/task"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

func Request(
	client *http.Client,
	request *http.Request,
	logger interfaces.LoggerInterface,
) (error, *http.Response, *string) {

	err, body, response := GetBody(client, request)

	if err != nil {
		color.Red("Request error: %v", err)
		return err, nil, nil
	}

	if len(*body) == 0 {
		return Errors.ResponseEmptyError, nil, nil
	}

	return nil, response, body
}

func NewRequest(task *task.TyphoonTask) (*http.Request, error) {
	color.Yellow("Create request %s : %s", task.GetFetcherMethod(), task.GetFetcherUrl())
	return http.NewRequest(task.GetFetcherMethod(), task.GetFetcherUrl(), task.GetRequestBody())
}

func FetchData(task *task.TyphoonTask) (error, *string) {
	_, client := GetHttpClientTransport(task)
	var err error
	request, err := NewRequest(task)
	err, data, _ := GetBody(client, request)
	return err, data
}

func CreateProxyRequestPipeline(opts *forms.Options) *HttpRequestPipeline {

	return &HttpRequestPipeline{
		BasePipeline: &forms.BasePipeline{
			MetaInfo: &label.MetaInfo{
				Name: "http-request",
			},
			Options: opts,
			Middlewares: []interfaces.MiddlewareInterface{
				ConstructorPrepareRequestMiddleware(true),
				ConstructorRequestHeaderMiddleware(true),
				ConstructorProxySettingMiddleware(true),
				ConstructorProxyRequestSettingsMiddleware(true),
			},
		},

		Fn: func(context context.Context,
			task interfaces.TaskInterface,
			logger interfaces.LoggerInterface,
			client *http.Client,
			request *http.Request,
			transport *http.Transport) (error, context.Context) {

			err, response, data := Request(client, request, logger)
			if err != nil {
				return err, nil
			}
			context = NewResponseCtx(context, response)
			context = NewResponseDataCtx(context, data)

			return nil, context

		},
		Cn: func(err error,
			context context.Context,
			task interfaces.TaskInterface,
			logger interfaces.LoggerInterface) {

			// Block current proxy
			proxy := task.GetProxyAddress()
			logger.Error("block proxy: %s", proxy)
			urlSupported := fmt.Sprintf("%s/block?url=%s&proxy=%s&code=599",
				task.GetProxyServerUrl(),
				task.GetFetcherUrl(), proxy,
			)
			logger.Info("block proxy :", urlSupported)

			errBlockRequest := retry.Do(func() error {
				client := GetHttpClient(task)
				request, errR := http.NewRequest(http.MethodGet, urlSupported, nil)
				if errR != nil {
					return errR
				}
				errR, body, _ := GetBody(client, request)
				if errR != nil || body == nil {
					return Errors.ProxyServerError
				}

				var proxyResponse models.Proxy
				err = utils.JsonLoad(&proxyResponse, *body)
				if err != nil {
					color.Red("JsonLoad has Error: %s", err.Error())
					return err
				}
				if !proxyResponse.Success {
					return Errors.ResponseNotOkError

				}

				color.Green(fmt.Sprintf("proxy %s was be blocked ", proxy))

				return nil

			})

			if errBlockRequest != nil {
				logger.Error("Fatal exception. Impossible block proxy.")
				os.Exit(1)
			}

		},
	}
}

func CreateRequestPipeline() *HttpRequestPipeline {
	return &HttpRequestPipeline{
		BasePipeline: &forms.BasePipeline{
			MetaInfo: &label.MetaInfo{
				Name: "http-request",
			},
			Middlewares: []interfaces.MiddlewareInterface{
				ConstructorPrepareRequestMiddleware(true),
				ConstructorRequestHeaderMiddleware(true),
			},
		},

		Fn: func(context context.Context,
			task interfaces.TaskInterface,
			logger interfaces.LoggerInterface,
			client *http.Client,
			request *http.Request,
			transport *http.Transport) (error, context.Context) {

			err, response, data := Request(client, request, logger)
			if err != nil {
				return err, nil
			}
			context = NewResponseCtx(context, response)
			context = NewResponseDataCtx(context, data)

			return nil, context

		},
		Cn: func(err error,
			context context.Context,
			task interfaces.TaskInterface,
			logger interfaces.LoggerInterface) {

			logger.Error("--- ", err.Error())
		},
	}
}
