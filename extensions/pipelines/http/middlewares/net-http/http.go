package net_http

import (
	"context"
	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/elements/forms"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/task"
	"net/http"
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

	taskInstance, okT := ctx.GetContextValue(context, TASK).(*task.TyphoonTask)
	request, okR := ctx.GetContextValue(context, REQUEST).(*http.Request)
	if !okT || !okR { reject(Errors.MiddlewareContextFailed); return }
	m.Fn(context, taskInstance, request, logger, reject, next)
}