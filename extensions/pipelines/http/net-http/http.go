package net_http

import (
	"context"
	"net/http"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/task"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

type HTTPRequestCallback func(
	context context.Context,
	task *task.TyphoonTask,
	request *http.Request,

	logger interfaces.LoggerInterface,
	reject func(err error),
	next func(ctx context.Context),
)


type HttpMiddleware struct {
	*forms.Middleware
	Fn HTTPRequestCallback
}

func (m *HttpMiddleware) Pass(
	context context.Context,
	logger interfaces.LoggerInterface,
	reject func(err error),
	next func(ctx context.Context),
	) {

	okT, taskInstance := task.Get(context)
	okR, request := GetRequestCtx(context)

	if !okT || !okR { reject(Errors.MiddlewareContextFailed); return }
	m.Fn(context, taskInstance, request, logger, reject, next)
}
