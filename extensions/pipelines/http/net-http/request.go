package net_http

import (
	"context"
	"net/http"

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

	if err != nil { return Errors.ResponseReadError, nil, nil }

	if len(*body) == 0 { return Errors.ResponseEmptyError, nil, nil }

	return nil, response, body
}

func NewRequest(task *task.TyphoonTask) (*http.Request, error)  {
	return http.NewRequest(task.GetFetcherMethod(), task.GetFetcherUrl(), nil)
}

func FetchData(task *task.TyphoonTask) (error, *string) {
	_, client := GetHttpClientTransport(task)
	var err error
	request, err := NewRequest(task)
	err, data, _ := GetBody(client, request)
	return err, data
}

func CreateRequestPipeline() *HttpRequestPipeline {
	return &HttpRequestPipeline{
		BasePipeline: &forms.BasePipeline{
			MetaInfo: &label.MetaInfo{
				Name: "http-request",
			},
		},
		Fn: func(context context.Context,
			task interfaces.TaskInterface,
			logger interfaces.LoggerInterface,
			client *http.Client,
			request *http.Request,
			transport *http.Transport) (error, context.Context) {

			err, response, data := Request(client, request, logger)
			if err != nil { return err, nil }
			context = NewResponseCtx(context, response)
			context = NewResponseDataCtx(context, data)

			return nil, context

		},
		Cn: func(err error,
			context context.Context,
			task interfaces.TaskInterface,
			logger interfaces.LoggerInterface) {

			logger.Error("TEST --- ",err.Error())
		},
	}
}