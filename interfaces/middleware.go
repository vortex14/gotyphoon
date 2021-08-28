package interfaces

import (
	"context"
	"github.com/gin-gonic/gin"
)

type ConstructorMiddleware func (required bool) MiddlewareInterface

type MiddlewareCallback func(context context.Context, loggerInterface LoggerInterface, reject func(err error))

type MiddlewareInterface interface {
	IsRequired() bool
	Pass(context context.Context, loggerInterface LoggerInterface, reject func(err error))

	MetaDataInterface
}

type Middleware struct {
	Name        string
	Required    bool
	Description string
	ctx         *gin.Context
	Callback    MiddlewareCallback
	PyCallback  MiddlewareCallback
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

func (m *Middleware) Pass(context context.Context, logger LoggerInterface, reject func(err error)) {
	m.Callback(context, logger, reject)
}