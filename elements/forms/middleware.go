package forms

import (
	"context"

	"github.com/vortex14/gotyphoon/interfaces"
)

type Middleware struct {
	Name        string
	Required    bool
	Description string
	Fn          interfaces.MiddlewareCallback
	PyCallback  interfaces.MiddlewareCallback
}

func (m *Middleware) GetName() string {
	return m.Name
}

func (m *Middleware) GetDescription() string {
	return m.Description
}

func (m *Middleware) IsRequired() bool {
	return m.Required
}

func (m *Middleware) Pass(
	context context.Context, logger interfaces.LoggerInterface,
	reject func(err error), next func(ctx context.Context),
	) {
	m.Fn(context, logger, reject, next)
}