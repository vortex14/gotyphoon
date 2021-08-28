package net_http

import (
	"context"
	"github.com/sirupsen/logrus"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
	"net/http"
)

const (
	BASICAuthUserTaskKey = "login"
	BASICAuthPasswordTaskKey = "password"

	NAMEHttpBasicAuthMiddleware = "Http Basic auth middleware"
	DESCRIPTIONHttpBasicAuthMiddleware = "Setting basic auth credentials for http request from Typhoon task"
)

type HTTPBasicAuthMiddleware struct {
	*interfaces.BaseLabel
}


func (m *HTTPBasicAuthMiddleware) Run(
	task *task.TyphoonTask,
	request *http.Request,
	) error {

	login, okL := task.Fetcher.Auth[BASICAuthUserTaskKey]

	passwd, okP := task.Fetcher.Auth[BASICAuthPasswordTaskKey]

	if okL && okP {
		request.SetBasicAuth(
			login,
			passwd,
		)
		logrus.Debug("setting are set for basic auth request")
		return nil
	}


	return Errors.MiddlewareBasicAuthOptionsNotFound
}


func (m *HTTPBasicAuthMiddleware) Pass(
	context context.Context,
	loggerInterface interfaces.LoggerInterface,
	reject func(err error),

	) {
	taskInstance, _ := context.Value(TASK).(*task.TyphoonTask)
	request, _ := context.Value(REQUEST).(*http.Request)
	if err := m.Run(taskInstance, request); err != nil {
		reject(err)
	}
}

func ConstructorBasicAuthMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HTTPBasicAuthMiddleware{
		BaseLabel: &interfaces.BaseLabel{
			Required: required,
			Name:        NAMEHttpBasicAuthMiddleware,
			Description: DESCRIPTIONHttpBasicAuthMiddleware,
		},
	}
}

