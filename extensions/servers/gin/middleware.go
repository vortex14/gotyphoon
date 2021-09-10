package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	
	"github.com/vortex14/gotyphoon/elements/forms"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
)

type RequestMiddlewareCallback func(
	context context.Context,
	ginCtx *gin.Context,

	logger interfaces.LoggerInterface,
	reject func(err error),
	next func(ctx context.Context),
)


type RequestMiddleware struct {
	*forms.Middleware
	Fn RequestMiddlewareCallback
}

func (m *RequestMiddleware) Pass(
	context context.Context,
	logger interfaces.LoggerInterface,
	reject func(err error),
	next func(ctx context.Context),
) {

	okT, ginCtx := GetRequestCtx(context)

	if !okT { reject(Errors.MiddlewareContextFailed); return }
	m.Fn(context, ginCtx, logger, reject, next)
}

