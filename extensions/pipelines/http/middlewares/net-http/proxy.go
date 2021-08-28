package net_http

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"

	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
)

type HTTPDefaultSetProxyMiddleware struct {
	*interfaces.BaseLabel
}

func (m *HTTPDefaultSetProxyMiddleware) Run(
	task *task.TyphoonTask,
	transport *http.Transport,
	) error {

	if task.Fetcher.IsProxyRequired == true {
		logrus.Debug("init proxy address ")
		proxyURL, err := url.Parse(task.Fetcher.Proxy)
		if err != nil {
			logrus.Error(err.Error())
			return Errors.ProxyUrlWrong
		}


		if proxyURL.Host != "" && proxyURL.Port() != "" {
			transport.Proxy = http.ProxyURL(proxyURL)
			logrus.Debug("task proxy ", proxyURL.Path)
		} else if proxyURL.Host == "" || proxyURL.Port() == "" {
			return Errors.ProxyTaskNotFound
		}
	}

	return nil
}

func (m *HTTPDefaultSetProxyMiddleware) Pass(
	context context.Context,
	loggerInterface interfaces.LoggerInterface,
	reject func(err error),
	) {

	taskInstance, _ := context.Value(TASK).(*task.TyphoonTask)
	transport, _ := context.Value(TRANSPORT).(*http.Transport)
	if err := m.Run(taskInstance, transport); err != nil {
		reject(err)
	}
}

func ConstructorProxyMiddleware(required bool) interfaces.MiddlewareInterface {
	return &HTTPDefaultSetProxyMiddleware{
		BaseLabel: &interfaces.BaseLabel{
			Required:    required,
			Name:        NAMEHttpBasicAuthMiddleware,
			Description: DESCRIPTIONHttpBasicAuthMiddleware,
		},
	}
}
