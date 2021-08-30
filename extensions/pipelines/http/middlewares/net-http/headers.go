package net_http

import (
	"context"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
	"net/http"
)

const (
	NAME = "Request Header Middleware"
	DESCRIPTION = "Setting request header for Typhoon task"
)

type HttpRequestSetHeaderMiddleware struct {
	*interfaces.BaseLabel
}

func (h *HttpRequestSetHeaderMiddleware) Run(
	task *task.TyphoonTask,
	request *http.Request,
	) error {

	for key, element := range task.Fetcher.Headers {
		request.Header.Add(
			key,
			element,
		)
	}
	return nil

}

func (h *HttpRequestSetHeaderMiddleware) Pass(
	context context.Context,
	loggerInterface interfaces.LoggerInterface,
	reject func(err error),

	) {
	task, _ := context.Value(TASK).(*task.TyphoonTask)
	request, _ := context.Value(REQUEST).(*http.Request)
	if err := h.Run(task, request); err != nil {
		reject(err)
	}
}

func ConstructorRequestHeaderMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HttpRequestSetHeaderMiddleware{
		BaseLabel: &interfaces.BaseLabel{
			Required: required,
			Name:        NAME,
			Description: DESCRIPTION,
		},
	}
}
