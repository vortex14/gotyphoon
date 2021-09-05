package forms

import (
	"context"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Middleware struct {
	*label.MetaInfo
	Fn          interfaces.MiddlewareCallback
	PyFn        interfaces.MiddlewareCallback
}

func (m *Middleware) Pass(
	context context.Context, logger interfaces.LoggerInterface,
	reject func(err error), next func(ctx context.Context),
	) {
	m.Fn(context, logger, reject, next)
}